package runner

import (
	"context"
	"testing"

	"breaklab/internal/progress"
)

func TestSetHypothesisRequiresActiveLab(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	if err := r.SetHypothesis(context.Background(), "a guess"); err != ErrNoActiveLab {
		t.Fatalf("err = %v, want ErrNoActiveLab", err)
	}
}

func TestSetHypothesisRecordsOnActiveAttempt(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	if err := r.SetHypothesis(ctx, "probably a config typo"); err != nil {
		t.Fatal(err)
	}

	s, _ := LoadSession(r.SessionPath)
	pf, _ := progress.Load(r.ProgressPath)
	if got := pf.Attempts[s.AttemptIndex].Hypothesis; got != "probably a config typo" {
		t.Errorf("hypothesis = %q", got)
	}
}
