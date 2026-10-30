package progress

import (
	"fmt"
	"sync"
	"time"
)

// mu serializes every read-modify-write in this file: the web server
// handles requests concurrently, and two overlapping writers both
// loading progress.json before either saves would silently drop
// whichever write loses the race.
var mu sync.Mutex

// AppendActive appends a new "active" attempt for lab/variant, started
// now, and returns its index into Attempts.
func AppendActive(path, lab, variant string) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	pf, err := Load(path)
	if err != nil {
		return 0, err
	}
	if pf.Version == 0 {
		pf.Version = 1
	}

	index := len(pf.Attempts)
	pf.Attempts = append(pf.Attempts, Attempt{
		ID:        newAttemptID(),
		Lab:       lab,
		Variant:   variant,
		StartedAt: time.Now().UTC(),
		Status:    StatusActive,
	})
	if err := Save(path, pf); err != nil {
		return 0, err
	}
	return index, nil
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
	mu.Lock()
	defer mu.Unlock()

	pf, err := Load(path)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(pf.Attempts) {
		return fmt.Errorf("progress.json: attempt index %d out of range (have %d attempts)", index, len(pf.Attempts))
	}
	pf.Attempts[index].Postmortem = postmortemPath
	return Save(path, pf)
}

// MarkComplete transitions the attempt at index from "solved" to
// "complete". It no-ops if the attempt isn't currently "solved" (not
// yet solved, already complete, or abandoned), so re-validating an
// already-complete postmortem is a safe no-op rather than an error.
func MarkComplete(path string, index int) error {
	mu.Lock()
	defer mu.Unlock()

	pf, err := Load(path)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(pf.Attempts) {
		return fmt.Errorf("progress.json: attempt index %d out of range (have %d attempts)", index, len(pf.Attempts))
	}
	if pf.Attempts[index].Status != StatusSolved {
		return nil
	}
	pf.Attempts[index].Status = StatusComplete
	return Save(path, pf)
}

// SkipStage appends one "skipped" attempt for every id in labIDs that has
// no attempt yet in progress.json — the roadmap's per-stage "I already
// know this" action. A lab with any existing attempt (active, solved,
// complete, or even abandoned) is left alone, so skip never overwrites
// real progress.
func SkipStage(path string, labIDs []string) error {
	mu.Lock()
	defer mu.Unlock()

	pf, err := Load(path)
	if err != nil {
		return err
	}
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
	return Save(path, pf)
}

// BackfillAttemptIDs assigns a new ID to every attempt in progress.json
// that doesn't have one yet — attempts recorded before Attempt.ID
// existed. `breaklab sync` calls this once at startup so every attempt
// is identifiable before push/pull; once every attempt already has an
// ID, this is a no-op that doesn't even rewrite the file.
func BackfillAttemptIDs(path string) error {
	mu.Lock()
	defer mu.Unlock()

	pf, err := Load(path)
	if err != nil {
		return err
	}

	changed := false
	for i := range pf.Attempts {
		if pf.Attempts[i].ID == "" {
			pf.Attempts[i].ID = newAttemptID()
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return Save(path, pf)
}

// update loads progress.json and applies fn to the attempt at index, but
// only while it is still "active" — every write in this file is
// idempotent against an attempt that has already been finalized
// (design.md Decision 6), so a repeat call is a safe no-op rather than
// corrupting a terminal record.
func update(path string, index int, fn func(*Attempt)) error {
	mu.Lock()
	defer mu.Unlock()

	pf, err := Load(path)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(pf.Attempts) {
		return fmt.Errorf("progress.json: attempt index %d out of range (have %d attempts)", index, len(pf.Attempts))
	}
	if pf.Attempts[index].Status != StatusActive {
		return nil
	}
	fn(&pf.Attempts[index])
	return Save(path, pf)
}
