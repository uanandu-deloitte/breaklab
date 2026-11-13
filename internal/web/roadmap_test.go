package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"breaklab/internal/progress"
)

func solvedAt(t time.Time) *time.Time { return &t }

// add-curriculum-phases task 6.2: the flat roadmap page is superseded
// by per-phase pages; its old URL redirects rather than 404ing.
func TestRoadmapPageRedirectsToPhases(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()

	req := httptest.NewRequest("GET", "/roadmap", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
	}
	if loc := rec.Header().Get("Location"); loc != "/phases" {
		t.Fatalf("Location = %q, want /phases", loc)
	}
}

func TestStageSkipAPIHandler(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()

	req := httptest.NewRequest("POST", "/api/phases/stages/linux/skip", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204 (body: %s)", rec.Code, rec.Body)
	}

	pf, err := progress.Load(cfg.ProgressPath)
	if err != nil {
		t.Fatalf("loading progress.json: %v", err)
	}
	if len(pf.Attempts) != 1 || pf.Attempts[0].Lab != "alpha" || pf.Attempts[0].Status != progress.StatusSkipped {
		t.Fatalf("attempts = %+v, want one skipped attempt for lab alpha", pf.Attempts)
	}
	if pf.Attempts[0].SolvedAt == nil {
		t.Error("skipped attempt should set SolvedAt")
	}
}

func TestStageSkipAPIHandlerLeavesExistingAttemptsAlone(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()
	if err := progress.Save(cfg.ProgressPath, &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "alpha", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/phases/stages/linux/skip", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204 (body: %s)", rec.Code, rec.Body)
	}

	pf, err := progress.Load(cfg.ProgressPath)
	if err != nil {
		t.Fatalf("loading progress.json: %v", err)
	}
	if len(pf.Attempts) != 1 || pf.Attempts[0].Status != progress.StatusComplete {
		t.Fatalf("attempts = %+v, want the original complete attempt untouched, no skipped one added", pf.Attempts)
	}
}
