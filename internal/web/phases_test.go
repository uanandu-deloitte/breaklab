package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"breaklab/internal/web/phases"
)

func TestPhaseListAPIEndToEnd(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/phases", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body)
	}

	var view phases.PhaseListView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if len(view.Phases) != 1 || view.Phases[0].ID != "phase-1" {
		t.Fatalf("phases = %+v, want exactly phase-1 (actionRoadmapYAML's fixture)", view.Phases)
	}
}

func TestPhaseDetailAPIEndToEnd(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/phases/phase-1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body)
	}

	var view phases.PhaseDetailView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if len(view.Stages) != 1 || view.Stages[0].ID != "linux" {
		t.Fatalf("stages = %+v, want exactly the linux stage", view.Stages)
	}
}

func TestPhaseDetailAPIUnknownPhaseIs404(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/phases/phase-nope", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 (body: %s)", rec.Code, rec.Body)
	}
}
