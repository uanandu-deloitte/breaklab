package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"breaklab/internal/progress"
)

func saveProgress(t *testing.T, path string, attempts ...progress.Attempt) {
	t.Helper()
	if err := progress.Save(path, &progress.ProgressFile{Version: 1, Attempts: attempts}); err != nil {
		t.Fatal(err)
	}
}

func TestPostmortemRejectsTraversalVariant(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	solved := time.Now()
	saveProgress(t, r.ProgressPath, progress.Attempt{
		Lab: "alpha", Variant: "../../x", StartedAt: solved.Add(-time.Minute),
		SolvedAt: &solved, Status: progress.StatusSolved,
	})
	root := filepath.Dir(r.LabsDir)
	before := countFiles(t, root)

	if _, err := r.Postmortem(context.Background()); err == nil {
		t.Fatal("expected an error for a traversal-shaped variant")
	}
	if after := countFiles(t, root); after != before {
		t.Errorf("files under the test root went from %d to %d; nothing should be created", before, after)
	}
}

func TestValidatePostmortemNoPathRejectsPathOutsideDir(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	outside := filepath.Join(filepath.Dir(r.LabsDir), "outside.md")
	if err := os.WriteFile(outside, []byte("---\nlab: alpha\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	solved := time.Now()
	saveProgress(t, r.ProgressPath, progress.Attempt{
		Lab: "alpha", Variant: "first-variant", StartedAt: solved.Add(-time.Minute),
		SolvedAt: &solved, Status: progress.StatusSolved, Postmortem: outside,
	})

	if _, err := r.ValidatePostmortem(context.Background(), ""); err == nil {
		t.Fatal("expected an error for a recorded path outside the postmortems dir")
	}
	pf, _ := progress.Load(r.ProgressPath)
	if pf.Attempts[0].Status != progress.StatusSolved {
		t.Errorf("status = %q, want unchanged solved", pf.Attempts[0].Status)
	}
}

func countFiles(t *testing.T, root string) int {
	t.Helper()
	n := 0
	_ = filepath.WalkDir(filepath.Dir(root), func(_ string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			n++
		}
		return nil
	})
	return n
}
