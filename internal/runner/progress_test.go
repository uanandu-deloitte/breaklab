package runner

import (
	"context"
	"os/exec"
	"testing"

	"breaklab/internal/progress"
)

func loadAttempt(t *testing.T, path string, index int) progress.Attempt {
	t.Helper()
	pf, err := progress.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if index < 0 || index >= len(pf.Attempts) {
		t.Fatalf("attempt index %d out of range (have %d)", index, len(pf.Attempts))
	}
	return pf.Attempts[index]
}

func TestStartCreatesActiveAttempt(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	s, err := LoadSession(r.SessionPath)
	if err != nil {
		t.Fatal(err)
	}
	a := loadAttempt(t, r.ProgressPath, s.AttemptIndex)
	if a.Lab != "alpha" || a.Variant != "first-variant" || a.Status != progress.StatusActive {
		t.Fatalf("attempt = %+v", a)
	}
}

func TestSwitchingLabsAbandonsThePreviousAttempt(t *testing.T) {
	r, _, _ := setup(t, "alpha", "beta")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)
	alphaIndex := s.AttemptIndex

	if err := r.Start(ctx, "beta", true); err != nil {
		t.Fatal(err)
	}

	a := loadAttempt(t, r.ProgressPath, alphaIndex)
	if a.Status != progress.StatusAbandoned {
		t.Errorf("alpha's attempt = %+v, want abandoned", a)
	}
}

func TestRestartingSameLabWithoutResetAbandonsThePreviousAttempt(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)
	firstIndex := s.AttemptIndex

	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s2, _ := LoadSession(r.SessionPath)
	if s2.AttemptIndex == firstIndex {
		t.Fatalf("restart reused the same attempt index: %d", s2.AttemptIndex)
	}

	a := loadAttempt(t, r.ProgressPath, firstIndex)
	if a.Status != progress.StatusAbandoned {
		t.Errorf("first attempt = %+v, want abandoned", a)
	}
}

func TestStopAbandonsAnActiveAttempt(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)
	index := s.AttemptIndex

	if err := r.Stop(ctx); err != nil {
		t.Fatal(err)
	}

	a := loadAttempt(t, r.ProgressPath, index)
	if a.Status != progress.StatusAbandoned {
		t.Errorf("attempt = %+v, want abandoned", a)
	}
}

func TestStopAfterSolvingLeavesAttemptSolved(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)
	index := s.AttemptIndex

	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}
	if err := r.Stop(ctx); err != nil {
		t.Fatal(err)
	}

	a := loadAttempt(t, r.ProgressPath, index)
	if a.Status != progress.StatusSolved {
		t.Errorf("attempt = %+v, want still solved (not overwritten by stop)", a)
	}
}

func TestCheckPassingMarksSolvedWithHintsSnapshot(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)
	if _, err := r.Hint(ctx); err != nil {
		t.Fatal(err)
	}

	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}

	a := loadAttempt(t, r.ProgressPath, s.AttemptIndex)
	if a.Status != progress.StatusSolved || a.SolvedAt == nil || a.HintsUsed != 1 {
		t.Fatalf("attempt = %+v", a)
	}
}

func TestCheckFailingIncrementsFailedChecks(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	failErr := exec.Command("bash", "-c", "exit 1").Run()
	r.HostExec = func(context.Context, string, string, []string, []string) ([]byte, error) {
		return []byte("FAIL: nope\n"), failErr
	}
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)

	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}

	a := loadAttempt(t, r.ProgressPath, s.AttemptIndex)
	if a.Status != progress.StatusActive || a.FailedChecks != 2 {
		t.Fatalf("attempt = %+v", a)
	}
}

func TestCheckAfterSolvedDoesNotReincrementOrChangeSolvedAt(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)

	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}
	first := loadAttempt(t, r.ProgressPath, s.AttemptIndex)

	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}
	second := loadAttempt(t, r.ProgressPath, s.AttemptIndex)

	if !second.SolvedAt.Equal(*first.SolvedAt) || second.FailedChecks != first.FailedChecks {
		t.Fatalf("re-checking after solved changed the record: before %+v, after %+v", first, second)
	}
}

func TestResetDoesNotModifyProgress(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)
	index := s.AttemptIndex

	before, err := progress.Load(r.ProgressPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Reset(ctx); err != nil {
		t.Fatal(err)
	}

	after, err := progress.Load(r.ProgressPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Attempts) != len(before.Attempts) {
		t.Fatalf("reset changed attempt count: before %d, after %d", len(before.Attempts), len(after.Attempts))
	}
	if after.Attempts[index] != before.Attempts[index] {
		t.Errorf("reset modified the attempt record: before %+v, after %+v", before.Attempts[index], after.Attempts[index])
	}

	s2, _ := LoadSession(r.SessionPath)
	if s2.AttemptIndex != index {
		t.Errorf("reset changed AttemptIndex: before %d, after %d", index, s2.AttemptIndex)
	}
}
