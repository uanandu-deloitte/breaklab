package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/web/capstone"
	"breaklab/internal/web/webtest"
)

const capstoneRoadmapYAML = `stages:
  - id: linux
    title: Linux fundamentals
  - id: finale
    title: Capstone
    capstone: true
`

// capstoneFixtureConfig lays out one regular-stage lab ("web-down") and
// one capstone-stage lab ("final-incident"), so buildCapstone has both
// a lockable regular stage and a real lab to render once unlocked.
func capstoneFixtureConfig(t *testing.T) Config {
	t.Helper()
	dir := t.TempDir()

	webtest.MustWriteFile(t, filepath.Join(dir, "roadmap.yaml"), capstoneRoadmapYAML)
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "web-down", "lab.yaml"), `id: web-down
title: The website returns nothing
stage: linux
difficulty: 1
estimated_minutes: 20
environment:
  dockerfile: Dockerfile
  ports: []
check:
  script: check.sh
  run_in: host
`)
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "web-down", "TICKET.md"), "# INC-1\n")
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "final-incident", "lab.yaml"), `id: final-incident
title: Everything is on fire
stage: finale
difficulty: 3
estimated_minutes: 45
environment:
  dockerfile: Dockerfile
  ports: []
check:
  script: check.sh
  run_in: host
`)
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "final-incident", "TICKET.md"), "# INC-2\n")
	webtest.MustWriteFile(t, filepath.Join(dir, "progress", "progress.json"), `{"version": 1, "attempts": []}`)

	return Config{
		LabsDir:      filepath.Join(dir, "labs"),
		RoadmapPath:  filepath.Join(dir, "roadmap.yaml"),
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		SessionPath:  filepath.Join(dir, ".breaklab", "session.json"),
		TemplatesDir: t.TempDir(),
		FrontendDir:  t.TempDir(),
	}
}

func TestCapstoneAPILockedBeforeRegularStagesComplete(t *testing.T) {
	cfg := capstoneFixtureConfig(t)

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/capstone", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var view capstone.CapstoneView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.Unlocked {
		t.Error("Unlocked = true before the regular stage is complete")
	}
	if len(view.Progress) != 1 || view.Progress[0].ID != "linux" || view.Progress[0].Completed != 0 || view.Progress[0].Total != 1 {
		t.Errorf("Progress = %+v, want one linux stage at 0/1", view.Progress)
	}
	if view.Lab.ID != "" {
		t.Errorf("Lab = %+v, want empty while locked", view.Lab)
	}
}

func TestCapstoneAPIUnlockedServesTheLab(t *testing.T) {
	cfg := capstoneFixtureConfig(t)
	solved := time.Now().UTC()
	pf := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "web-down", StartedAt: solved.Add(-time.Minute), SolvedAt: &solved, Status: progress.StatusComplete},
	}}
	if err := progress.Save(cfg.ProgressPath, pf); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/capstone", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var view capstone.CapstoneView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if !view.Unlocked {
		t.Fatal("Unlocked = false after the only regular stage is complete")
	}
	if view.Lab.ID != "final-incident" {
		t.Errorf("Lab.ID = %q, want final-incident", view.Lab.ID)
	}
}

func TestCapstoneAPINotFoundWithNoCapstoneStage(t *testing.T) {
	cfg := capstoneFixtureConfig(t)
	webtest.MustWriteFile(t, cfg.RoadmapPath, "stages:\n  - id: linux\n    title: Linux fundamentals\n")

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/capstone", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 (body: %s)", rec.Code, rec.Body.String())
	}
}
