package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"breaklab/internal/progress"
)

func TestStatusNoActiveLab(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	view, err := r.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if view.Active != nil {
		t.Errorf("Active = %+v, want nil", view.Active)
	}
	if len(view.Stages) != 1 || view.Stages[0].Total != 1 || view.Stages[0].Completed != 0 {
		t.Errorf("Stages = %+v", view.Stages)
	}
}

func TestStatusActiveLabElapsedAndHints(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Hint(ctx); err != nil {
		t.Fatal(err)
	}

	view, err := r.Status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if view.Active == nil {
		t.Fatal("Active = nil, want a block")
	}
	if view.Active.Lab != "alpha" || view.Active.HintsUsed != 1 || view.Active.ElapsedMinutes < 0 {
		t.Errorf("Active = %+v", view.Active)
	}
}

func TestStatusStageCompletionAfterSolve(t *testing.T) {
	r, _, _ := setup(t, "alpha", "beta")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}

	view, err := r.Status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Stages) != 1 || view.Stages[0].Total != 2 || view.Stages[0].Completed != 1 {
		t.Fatalf("Stages = %+v", view.Stages)
	}
	if view.Streak != 1 {
		t.Errorf("Streak = %d, want 1", view.Streak)
	}
}

func TestStatusListsReviewDueLabs(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	old := time.Now().Add(-(progress.ReviewInterval + time.Hour))
	pf := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "alpha", Variant: "first-variant", SolvedAt: &old, Status: progress.StatusComplete},
	}}
	if err := progress.Save(r.ProgressPath, pf); err != nil {
		t.Fatal(err)
	}

	view, err := r.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(view.ReviewDue) != 1 || view.ReviewDue[0] != "alpha" {
		t.Errorf("ReviewDue = %v, want [alpha]", view.ReviewDue)
	}
}

// appendLabTags appends a "tags: [...]" line to id's generated
// lab.yaml — setup()'s fixture labs carry none by default, and this is
// simpler than hand-rolling a second lab.yaml template just for this
// test's tag data.
func appendLabTags(t *testing.T, labsDir, id string, tags ...string) {
	t.Helper()
	f, err := os.OpenFile(filepath.Join(labsDir, id, "lab.yaml"), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("tags: [" + strings.Join(tags, ", ") + "]\n"); err != nil {
		t.Fatal(err)
	}
}

func TestStatusListsPatterns(t *testing.T) {
	r, _, _ := setup(t, "alpha", "beta")
	appendLabTags(t, r.LabsDir, "alpha", "networking")
	appendLabTags(t, r.LabsDir, "beta", "linux")

	pf := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "alpha", Status: progress.StatusComplete, HintsUsed: 6},
		{Lab: "alpha", Status: progress.StatusComplete, HintsUsed: 6},
		{Lab: "alpha", Status: progress.StatusComplete, HintsUsed: 6},
		{Lab: "beta", Status: progress.StatusComplete, HintsUsed: 1},
		{Lab: "beta", Status: progress.StatusComplete, HintsUsed: 1},
		{Lab: "beta", Status: progress.StatusComplete, HintsUsed: 1},
	}}
	if err := progress.Save(r.ProgressPath, pf); err != nil {
		t.Fatal(err)
	}

	view, err := r.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Patterns) != 1 || view.Patterns[0].Tag != "networking" {
		t.Fatalf("Patterns = %+v, want exactly one pattern for networking", view.Patterns)
	}
}

func TestStatusNoPatternsWhenNoneQualify(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	appendLabTags(t, r.LabsDir, "alpha", "networking")
	pf := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "alpha", Status: progress.StatusComplete, HintsUsed: 1},
	}}
	if err := progress.Save(r.ProgressPath, pf); err != nil {
		t.Fatal(err)
	}

	view, err := r.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Patterns) != 0 {
		t.Errorf("Patterns = %+v, want empty: only 1 sample, below the floor", view.Patterns)
	}
}

func TestStatusSessionWithUnknownAttemptIndexErrors(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, _ := LoadSession(r.SessionPath)
	s.AttemptIndex = 99
	if err := SaveSession(r.SessionPath, *s); err != nil {
		t.Fatal(err)
	}

	if _, err := r.Status(ctx); err == nil {
		t.Fatal("expected an error for an out-of-range attempt index")
	}
}
