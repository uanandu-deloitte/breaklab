package progress

import (
	"path/filepath"
	"testing"
)

func TestAppendActive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")

	i0, err := AppendActive(path, "web-down", "config-typo")
	if err != nil || i0 != 0 {
		t.Fatalf("i0 = %d, err = %v", i0, err)
	}
	i1, err := AppendActive(path, "web-down", "port-taken")
	if err != nil || i1 != 1 {
		t.Fatalf("i1 = %d, err = %v", i1, err)
	}

	pf, err := Load(path)
	if err != nil || len(pf.Attempts) != 2 {
		t.Fatalf("pf = %+v, err = %v", pf, err)
	}
	a := pf.Attempts[0]
	if a.Lab != "web-down" || a.Variant != "config-typo" || a.Status != StatusActive ||
		a.FailedChecks != 0 || a.HintsUsed != 0 || a.SolvedAt != nil {
		t.Fatalf("attempt 0 = %+v", a)
	}
	if a.ID == "" {
		t.Error("attempt 0 ID should be non-empty")
	}
	if pf.Attempts[1].ID == "" || pf.Attempts[1].ID == a.ID {
		t.Errorf("attempt 1 ID = %q, want non-empty and distinct from attempt 0's %q", pf.Attempts[1].ID, a.ID)
	}
}

func TestSkipStage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")

	if err := SkipStage(path, []string{"lnx-001", "lnx-002"}); err != nil {
		t.Fatal(err)
	}

	pf, err := Load(path)
	if err != nil || len(pf.Attempts) != 2 {
		t.Fatalf("pf = %+v, err = %v", pf, err)
	}
	for _, a := range pf.Attempts {
		if a.Status != StatusSkipped || a.SolvedAt == nil {
			t.Errorf("attempt %+v, want status skipped with SolvedAt set", a)
		}
	}
}

func TestSkipStageLeavesAttemptedLabsAlone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	if _, err := AppendActive(path, "lnx-001", ""); err != nil {
		t.Fatal(err)
	}

	if err := SkipStage(path, []string{"lnx-001", "lnx-002"}); err != nil {
		t.Fatal(err)
	}

	pf, err := Load(path)
	if err != nil || len(pf.Attempts) != 2 {
		t.Fatalf("pf = %+v, err = %v, want the original lnx-001 attempt plus one new skipped lnx-002", pf, err)
	}
	if pf.Attempts[0].Status != StatusActive {
		t.Errorf("lnx-001's existing attempt should be untouched, got status %q", pf.Attempts[0].Status)
	}
	if pf.Attempts[1].Lab != "lnx-002" || pf.Attempts[1].Status != StatusSkipped {
		t.Errorf("lnx-002 should get a new skipped attempt, got %+v", pf.Attempts[1])
	}
}

func TestBackfillAttemptIDsAssignsOnlyMissingOnes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	if err := Save(path, &ProgressFile{Version: 1, Attempts: []Attempt{
		{ID: "already-has-one", Lab: "web-down", Status: StatusComplete},
		{Lab: "stuck-lock", Status: StatusComplete},
	}}); err != nil {
		t.Fatal(err)
	}

	if err := BackfillAttemptIDs(path); err != nil {
		t.Fatal(err)
	}
	pf, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if pf.Attempts[0].ID != "already-has-one" {
		t.Errorf("existing ID should be untouched, got %q", pf.Attempts[0].ID)
	}
	backfilled := pf.Attempts[1].ID
	if backfilled == "" {
		t.Fatal("stuck-lock attempt should have been assigned an ID")
	}

	// Re-running once every attempt has an ID must be a true no-op.
	if err := BackfillAttemptIDs(path); err != nil {
		t.Fatal(err)
	}
	pf2, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if pf2.Attempts[1].ID != backfilled {
		t.Errorf("re-running backfill changed an already-assigned ID: got %q, want %q", pf2.Attempts[1].ID, backfilled)
	}
}

func TestMarkSolvedIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	index, err := AppendActive(path, "web-down", "config-typo")
	if err != nil {
		t.Fatal(err)
	}

	if err := MarkSolved(path, index, 1); err != nil {
		t.Fatal(err)
	}
	pf, _ := Load(path)
	first := pf.Attempts[index]
	if first.Status != StatusSolved || first.SolvedAt == nil || first.HintsUsed != 1 {
		t.Fatalf("after first solve: %+v", first)
	}

	// A second solve (e.g. re-running `breaklab check` after already
	// passing) must not change solved_at or hints_used.
	if err := MarkSolved(path, index, 99); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	second := pf.Attempts[index]
	if !second.SolvedAt.Equal(*first.SolvedAt) || second.HintsUsed != 1 {
		t.Fatalf("second solve changed the record: before %+v, after %+v", first, second)
	}
}

func TestMarkAbandonedOnlyOverwritesActive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	index, err := AppendActive(path, "web-down", "config-typo")
	if err != nil {
		t.Fatal(err)
	}

	if err := MarkAbandoned(path, index, 2); err != nil {
		t.Fatal(err)
	}
	pf, _ := Load(path)
	if a := pf.Attempts[index]; a.Status != StatusAbandoned || a.HintsUsed != 2 {
		t.Fatalf("abandoned = %+v", a)
	}

	// A fresh attempt that gets solved must not be flipped to abandoned
	// by a later call (e.g. Start finalizing a stale session reference).
	index2, err := AppendActive(path, "web-down", "port-taken")
	if err != nil {
		t.Fatal(err)
	}
	if err := MarkSolved(path, index2, 0); err != nil {
		t.Fatal(err)
	}
	if err := MarkAbandoned(path, index2, 5); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	if a := pf.Attempts[index2]; a.Status != StatusSolved {
		t.Fatalf("solved attempt was overwritten: %+v", a)
	}
}

func TestIncrementFailedChecksStopsAfterTerminal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	index, err := AppendActive(path, "web-down", "config-typo")
	if err != nil {
		t.Fatal(err)
	}

	if err := IncrementFailedChecks(path, index); err != nil {
		t.Fatal(err)
	}
	if err := IncrementFailedChecks(path, index); err != nil {
		t.Fatal(err)
	}
	pf, _ := Load(path)
	if a := pf.Attempts[index]; a.FailedChecks != 2 {
		t.Fatalf("failed_checks = %d, want 2", a.FailedChecks)
	}

	if err := MarkSolved(path, index, 0); err != nil {
		t.Fatal(err)
	}
	if err := IncrementFailedChecks(path, index); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	if a := pf.Attempts[index]; a.FailedChecks != 2 {
		t.Fatalf("failed_checks after solve = %d, want still 2", a.FailedChecks)
	}
}

func TestSetHypothesisOnlyWhileActive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	index, err := AppendActive(path, "web-down", "config-typo")
	if err != nil {
		t.Fatal(err)
	}

	if err := SetHypothesis(path, index, "nginx config has a syntax error"); err != nil {
		t.Fatal(err)
	}
	pf, _ := Load(path)
	if a := pf.Attempts[index]; a.Hypothesis != "nginx config has a syntax error" {
		t.Fatalf("hypothesis = %q", a.Hypothesis)
	}

	// Rewriting while still active replaces it.
	if err := SetHypothesis(path, index, "actually, the port is taken"); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	if a := pf.Attempts[index]; a.Hypothesis != "actually, the port is taken" {
		t.Fatalf("hypothesis after rewrite = %q", a.Hypothesis)
	}

	if err := MarkSolved(path, index, 0); err != nil {
		t.Fatal(err)
	}
	// Once solved, a further write is a no-op (gated on active, like
	// every other write in this file).
	if err := SetHypothesis(path, index, "ignored"); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	if a := pf.Attempts[index]; a.Hypothesis != "actually, the port is taken" {
		t.Fatalf("hypothesis changed after solve: %q", a.Hypothesis)
	}
}

func TestSetPostmortemIgnoresStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	index, err := AppendActive(path, "web-down", "config-typo")
	if err != nil {
		t.Fatal(err)
	}

	// Set on a still-active attempt (e.g. a hypothetical future use);
	// the write isn't gated on status the way the others are.
	if err := SetPostmortem(path, index, "postmortems/2026-10-17-web-down-config-typo.md"); err != nil {
		t.Fatal(err)
	}
	pf, _ := Load(path)
	if a := pf.Attempts[index]; a.Postmortem != "postmortems/2026-10-17-web-down-config-typo.md" || a.Status != StatusActive {
		t.Fatalf("attempt = %+v", a)
	}
}

func TestMarkCompletePromotesSolvedOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")

	activeIndex, err := AppendActive(path, "web-down", "config-typo")
	if err != nil {
		t.Fatal(err)
	}
	if err := MarkComplete(path, activeIndex); err != nil {
		t.Fatal(err)
	}
	pf, _ := Load(path)
	if a := pf.Attempts[activeIndex]; a.Status != StatusActive {
		t.Fatalf("still-active attempt was promoted: %+v", a)
	}

	abandonedIndex, err := AppendActive(path, "web-down", "port-taken")
	if err != nil {
		t.Fatal(err)
	}
	if err := MarkAbandoned(path, abandonedIndex, 0); err != nil {
		t.Fatal(err)
	}
	if err := MarkComplete(path, abandonedIndex); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	if a := pf.Attempts[abandonedIndex]; a.Status != StatusAbandoned {
		t.Fatalf("abandoned attempt was promoted: %+v", a)
	}

	solvedIndex, err := AppendActive(path, "web-down", "bad-perms")
	if err != nil {
		t.Fatal(err)
	}
	if err := MarkSolved(path, solvedIndex, 1); err != nil {
		t.Fatal(err)
	}
	if err := MarkComplete(path, solvedIndex); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	if a := pf.Attempts[solvedIndex]; a.Status != StatusComplete {
		t.Fatalf("solved attempt was not promoted: %+v", a)
	}

	// Already complete: a repeat call is a no-op, not an error.
	if err := MarkComplete(path, solvedIndex); err != nil {
		t.Fatal(err)
	}
	pf, _ = Load(path)
	if a := pf.Attempts[solvedIndex]; a.Status != StatusComplete {
		t.Fatalf("re-marking complete changed the record: %+v", a)
	}
}

func TestUpdateRejectsOutOfRangeIndex(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	if err := MarkSolved(path, 0, 0); err == nil {
		t.Fatal("expected an error for an out-of-range index on an empty file")
	}
}
