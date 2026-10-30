package runner

import (
	"context"
	"errors"
	"testing"
	"time"

	"breaklab/internal/progress"
)

func TestHintProgressesThroughGeneralThenVariantHints(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	// setup()'s fixture lab declares two general hints and one hint for
	// "first-variant" (the variant forced by setup()'s default env).
	want := []string{"general hint one", "general hint two", "first-variant hint"}
	for i, w := range want {
		hint, err := r.Hint(ctx)
		if err != nil || hint != w {
			t.Fatalf("hint %d: got %q, %v; want %q", i, hint, err, w)
		}
	}

	if _, err := r.Hint(ctx); !errors.Is(err, ErrOutOfHints) {
		t.Fatalf("after last hint: err = %v, want ErrOutOfHints", err)
	}
}

func TestHintNoActiveLab(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	if _, err := r.Hint(context.Background()); !errors.Is(err, ErrNoActiveLab) {
		t.Fatalf("err = %v, want ErrNoActiveLab", err)
	}
}

func TestHintMissingVariant(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	// Simulate lab.yaml losing the active variant between commands.
	s, err := LoadSession(r.SessionPath)
	if err != nil {
		t.Fatal(err)
	}
	s.Variant = "no-such-variant"
	if err := SaveSession(r.SessionPath, *s); err != nil {
		t.Fatal(err)
	}

	if _, err := r.Hint(ctx); err == nil {
		t.Fatal("expected an error for a missing variant")
	}
}

func TestResetZeroesHintsUsed(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	if _, err := r.Hint(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Hint(ctx); err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadSession(r.SessionPath); s.HintsUsed != 2 {
		t.Fatalf("HintsUsed before reset = %d, want 2", s.HintsUsed)
	}

	if err := r.Reset(ctx); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSession(r.SessionPath)
	if err != nil || s.HintsUsed != 0 {
		t.Fatalf("HintsUsed after reset = %+v, err %v, want 0", s, err)
	}

	hint, err := r.Hint(ctx)
	if err != nil || hint != "general hint one" {
		t.Fatalf("hint after reset = %q, err %v, want the lab's first hint", hint, err)
	}
}

func TestHintRefusedDuringReviewAttempt(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()

	// Seed a prior complete attempt for "alpha" before this session's
	// own attempt exists, so the attempt Start is about to create is a
	// review attempt.
	solved := time.Now().Add(-24 * time.Hour)
	pf := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "alpha", Variant: "first-variant", SolvedAt: &solved, Status: progress.StatusComplete},
	}}
	if err := progress.Save(r.ProgressPath, pf); err != nil {
		t.Fatal(err)
	}

	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	if _, err := r.Hint(ctx); !errors.Is(err, ErrReviewHintsUnavailable) {
		t.Fatalf("err = %v, want ErrReviewHintsUnavailable", err)
	}
}

func TestHintIncrementsAndPersistsAcrossCalls(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	if _, err := r.Hint(ctx); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSession(r.SessionPath)
	if err != nil || s.HintsUsed != 1 {
		t.Fatalf("session = %+v, err %v, want HintsUsed=1", s, err)
	}

	// A fresh Runner reading the same session file (as a new CLI process
	// would) continues from where the last call left off.
	r2 := &Runner{
		Docker:      r.Docker,
		LabsDir:     r.LabsDir,
		RoadmapPath: r.RoadmapPath,
		SessionPath: r.SessionPath,
		HistoryPath: r.HistoryPath,
	}
	hint, err := r2.Hint(ctx)
	if err != nil || hint != "general hint two" {
		t.Fatalf("hint = %q, err %v, want %q", hint, err, "general hint two")
	}
}
