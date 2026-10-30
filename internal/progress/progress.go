// Package progress keeps the learner's committed attempt log in
// progress/progress.json and the read-side summaries (`breaklab status`)
// derived from it.
package progress

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Attempt statuses. "complete" is declared for forward compatibility with
// postmortems (a later change promotes a "solved" attempt to "complete"
// once its postmortem is written) but nothing in this package ever
// produces it.
const (
	StatusActive    = "active"
	StatusSolved    = "solved"
	StatusComplete  = "complete"
	StatusAbandoned = "abandoned"

	// StatusSkipped marks a lab the learner declared already-known
	// (the roadmap's per-stage "skip" action) rather than solved. It
	// counts as done for lab/stage status and the capstone gate, but
	// never as a genuine solve: it's excluded from streaks, average
	// resolution time, and postmortem counts.
	StatusSkipped = "skipped"
)

// Attempt is one entry in progress.json.
type Attempt struct {
	// ID stably identifies this attempt across machines, so
	// `breaklab sync` can tell whether it has already pushed/pulled a
	// given attempt without relying on its (local-only, non-portable)
	// position in Attempts. Empty for attempts recorded before sync
	// existed; BackfillAttemptIDs assigns one lazily on first sync.
	ID           string     `json:"id,omitempty"`
	Lab          string     `json:"lab"`
	Variant      string     `json:"variant"`
	StartedAt    time.Time  `json:"started_at"`
	SolvedAt     *time.Time `json:"solved_at,omitempty"`
	FailedChecks int        `json:"failed_checks"`
	HintsUsed    int        `json:"hints_used"`
	Postmortem   string     `json:"postmortem,omitempty"`
	Status       string     `json:"status"`

	// Hypothesis is the learner's own guess at the cause, written while
	// the attempt is still active — before the real root cause is
	// known — so it can be compared against Root cause once the
	// postmortem is written. Optional; empty for an attempt where none
	// was recorded.
	Hypothesis string `json:"hypothesis,omitempty"`
}

// newAttemptID returns a random 128-bit hex string, unique enough to
// identify one attempt across every machine it's ever synced to or
// from.
func newAttemptID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err) // crypto/rand.Read only fails if the OS entropy source is broken
	}
	return hex.EncodeToString(b[:])
}

// ProgressFile is the parsed shape of progress/progress.json.
type ProgressFile struct {
	Version  int       `json:"version"`
	Attempts []Attempt `json:"attempts"`
}

// Load reads the progress file at path. A missing file is not an error —
// it returns an empty, version-1 ProgressFile, since progress.json is
// created on first use rather than shipped with the repo.
func Load(path string) (*ProgressFile, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &ProgressFile{Version: 1}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var pf ProgressFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("%s is corrupt (%v)", path, err)
	}
	return &pf, nil
}

// Save writes pf atomically (temp file + rename) with mode 0644: unlike
// .breaklab/session.json and .breaklab/history.json, progress.json is a
// committed artifact the learner is meant to read, edit, and `git add`,
// not machine-local run state.
func Save(path string, pf *ProgressFile) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "progress-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
