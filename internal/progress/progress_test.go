package progress

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress", "progress.json")
	pf, err := Load(path)
	if err != nil || pf.Version != 1 || len(pf.Attempts) != 0 {
		t.Fatalf("pf = %+v, err = %v", pf, err)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress", "progress.json")
	solvedAt := time.Date(2026, 10, 17, 10, 20, 0, 0, time.UTC)

	in := &ProgressFile{
		Version: 1,
		Attempts: []Attempt{
			{
				Lab:          "web-down",
				Variant:      "config-typo",
				StartedAt:    time.Date(2026, 10, 17, 10, 2, 0, 0, time.UTC),
				SolvedAt:     &solvedAt,
				FailedChecks: 2,
				HintsUsed:    1,
				Status:       StatusSolved,
			},
		},
	}
	if err := Save(path, in); err != nil {
		t.Fatal(err)
	}

	out, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Attempts) != 1 {
		t.Fatalf("attempts = %+v", out.Attempts)
	}
	got := out.Attempts[0]
	want := in.Attempts[0]
	if got.Lab != want.Lab || got.Variant != want.Variant || got.FailedChecks != want.FailedChecks ||
		got.HintsUsed != want.HintsUsed || got.Status != want.Status ||
		!got.StartedAt.Equal(want.StartedAt) || got.SolvedAt == nil || !got.SolvedAt.Equal(*want.SolvedAt) {
		t.Fatalf("round trip: got %+v, want %+v", got, want)
	}

	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o644 {
		t.Errorf("mode = %v, err %v, want 0644", info, err)
	}
}

func TestLoadFileWithoutIDStillParses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")
	// Pre-sync fixture: no "id" field on the attempt at all.
	body := `{"version":1,"attempts":[{"lab":"web-down","variant":"config-typo","started_at":"2026-10-17T10:02:00Z","failed_checks":0,"hints_used":0,"status":"active"}]}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	pf, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(pf.Attempts) != 1 || pf.Attempts[0].ID != "" || pf.Attempts[0].Lab != "web-down" {
		t.Fatalf("attempts = %+v", pf.Attempts)
	}
}

func TestLoadCorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("err = %v", err)
	}
}
