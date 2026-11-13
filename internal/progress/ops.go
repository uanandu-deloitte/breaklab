package progress

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// mu serializes every read-modify-write within this process: the web
// server handles requests concurrently, so two overlapping writers
// both loading progress.json before either saves would silently drop
// whichever write loses the race.
var mu sync.Mutex

// Package-level vars, not consts, so tests can shrink them instead of
// waiting out real multi-second timeouts.
var (
	lockRetryInterval = 20 * time.Millisecond
	lockTimeout       = 5 * time.Second
	staleLockAge      = 10 * time.Second
)

// acquireFileLock takes a cross-process advisory lock on path's
// sibling .lock file. mu above only ever serializes goroutines inside
// this one process — it does nothing for a second OS process, e.g.
// `breaklab serve` left running alongside a separate `breaklab start`
// invocation in another terminal, both racing to read-modify-write the
// same progress.json. It retries with a short backoff rather than
// blocking forever, and breaks a lock older than staleLockAge — left
// behind by a process that crashed mid-write — rather than deadlocking
// permanently.
func acquireFileLock(path string) (release func(), err error) {
	// progress.json's directory may not exist yet on the very first
	// write (Save creates it lazily) — create it here too, so locking
	// the very first write doesn't fail with ENOENT before Save ever
	// gets a chance to.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("progress.json: creating %s: %w", filepath.Dir(path), err)
	}

	lockPath := path + ".lock"
	deadline := time.Now().Add(lockTimeout)
	for {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			f.Close()
			return func() { os.Remove(lockPath) }, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("progress.json: creating lock file: %w", err)
		}
		if info, statErr := os.Stat(lockPath); statErr == nil && time.Since(info.ModTime()) > staleLockAge {
			os.Remove(lockPath)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("progress.json: timed out waiting for a lock held by another breaklab process (%s)", lockPath)
		}
		time.Sleep(lockRetryInterval)
	}
}

// withLock loads progress.json, lets fn read/mutate it, and saves the
// result — all under both mu (this process) and acquireFileLock
// (every process), so no concurrent writer anywhere can interleave
// with this read-modify-write cycle. fn returning save=false (e.g. an
// idempotent no-op guard, same as update's old behavior) skips the
// write entirely rather than rewriting an unchanged file.
func withLock(path string, fn func(pf *ProgressFile) (save bool, err error)) error {
	mu.Lock()
	defer mu.Unlock()

	release, err := acquireFileLock(path)
	if err != nil {
		return err
	}
	defer release()

	pf, err := Load(path)
	if err != nil {
		return err
	}
	save, err := fn(pf)
	if err != nil {
		return err
	}
	if !save {
		return nil
	}
	return Save(path, pf)
}

// AppendActive appends a new "active" attempt for lab/variant, started
// now, and returns its index into Attempts.
func AppendActive(path, lab, variant string) (int, error) {
	var index int
	err := withLock(path, func(pf *ProgressFile) (bool, error) {
		if pf.Version == 0 {
			pf.Version = 1
		}
		index = len(pf.Attempts)
		pf.Attempts = append(pf.Attempts, Attempt{
			ID:        newAttemptID(),
			Lab:       lab,
			Variant:   variant,
			StartedAt: time.Now().UTC(),
			Status:    StatusActive,
		})
		return true, nil
	})
	return index, err
}

// MarkSolved transitions the attempt at index to "solved", recording
// solved_at and hintsUsed. It no-ops if the attempt is no longer active.
func MarkSolved(path string, index, hintsUsed int) error {
	return update(path, index, func(a *Attempt) {
		now := time.Now().UTC()
		a.Status = StatusSolved
		a.SolvedAt = &now
		a.HintsUsed = hintsUsed
	})
}

// MarkAbandoned transitions the attempt at index to "abandoned",
// recording hintsUsed. It no-ops if the attempt is no longer active.
func MarkAbandoned(path string, index, hintsUsed int) error {
	return update(path, index, func(a *Attempt) {
		a.Status = StatusAbandoned
		a.HintsUsed = hintsUsed
	})
}

// IncrementFailedChecks bumps the attempt at index's failed_checks
// count. It no-ops if the attempt is no longer active.
func IncrementFailedChecks(path string, index int) error {
	return update(path, index, func(a *Attempt) {
		a.FailedChecks++
	})
}

// SetHypothesis records the learner's guess at the cause on the attempt
// at index, while it is still "active" — like every other write in this
// file gated through update, so a hypothesis can't be rewritten after
// the attempt has already been solved or abandoned.
func SetHypothesis(path string, index int, text string) error {
	return update(path, index, func(a *Attempt) {
		a.Hypothesis = text
	})
}

// SetPostmortem records postmortemPath on the attempt at index,
// regardless of its current status — unlike every other write in this
// file, this isn't gated on a specific prior status, since it only
// records where a file was created, not a lifecycle transition.
func SetPostmortem(path string, index int, postmortemPath string) error {
	return withLock(path, func(pf *ProgressFile) (bool, error) {
		if index < 0 || index >= len(pf.Attempts) {
			return false, fmt.Errorf("progress.json: attempt index %d out of range (have %d attempts)", index, len(pf.Attempts))
		}
		pf.Attempts[index].Postmortem = postmortemPath
		return true, nil
	})
}

// MarkComplete transitions the attempt at index from "solved" to
// "complete". It no-ops if the attempt isn't currently "solved" (not
// yet solved, already complete, or abandoned), so re-validating an
// already-complete postmortem is a safe no-op rather than an error.
func MarkComplete(path string, index int) error {
	return withLock(path, func(pf *ProgressFile) (bool, error) {
		if index < 0 || index >= len(pf.Attempts) {
			return false, fmt.Errorf("progress.json: attempt index %d out of range (have %d attempts)", index, len(pf.Attempts))
		}
		if pf.Attempts[index].Status != StatusSolved {
			return false, nil
		}
		pf.Attempts[index].Status = StatusComplete
		return true, nil
	})
}

// SkipStage appends one "skipped" attempt for every id in labIDs that has
// no attempt yet in progress.json — the roadmap's per-stage "I already
// know this" action. A lab with any existing attempt (active, solved,
// complete, or even abandoned) is left alone, so skip never overwrites
// real progress.
func SkipStage(path string, labIDs []string) error {
	return withLock(path, func(pf *ProgressFile) (bool, error) {
		if pf.Version == 0 {
			pf.Version = 1
		}

		attempted := make(map[string]bool, len(pf.Attempts))
		for _, a := range pf.Attempts {
			attempted[a.Lab] = true
		}

		now := time.Now().UTC()
		for _, id := range labIDs {
			if attempted[id] {
				continue
			}
			pf.Attempts = append(pf.Attempts, Attempt{
				ID:        newAttemptID(),
				Lab:       id,
				StartedAt: now,
				SolvedAt:  &now,
				Status:    StatusSkipped,
			})
		}
		return true, nil
	})
}

// update loads progress.json and applies fn to the attempt at index, but
// only while it is still "active" — every write in this file is
// idempotent against an attempt that has already been finalized
// (design.md Decision 6), so a repeat call is a safe no-op rather than
// corrupting a terminal record.
func update(path string, index int, fn func(*Attempt)) error {
	return withLock(path, func(pf *ProgressFile) (bool, error) {
		if index < 0 || index >= len(pf.Attempts) {
			return false, fmt.Errorf("progress.json: attempt index %d out of range (have %d attempts)", index, len(pf.Attempts))
		}
		if pf.Attempts[index].Status != StatusActive {
			return false, nil
		}
		fn(&pf.Attempts[index])
		return true, nil
	})
}
