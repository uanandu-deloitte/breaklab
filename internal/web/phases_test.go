package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/web/phases"
	"breaklab/internal/web/webcore"
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

// Mirrors TestRoadmapAPIAccountAwareReadPath: phase detail's lab status
// must be computed only from the local progress.json with no session,
// and from the account's merged attempts with a valid session.
func TestPhaseDetailAPIAccountAwareReadPath(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()
	cfg.AccountStore = newTestAccountStore(t)

	u, err := cfg.AccountStore.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	solved := time.Now().UTC()
	if err := cfg.AccountStore.PushAttempts(u.ID, []progress.Attempt{
		{ID: "hosted-1", Lab: "alpha", Status: progress.StatusComplete, StartedAt: solved.Add(-time.Minute), SolvedAt: &solved},
	}); err != nil {
		t.Fatal(err)
	}
	token, _, err := cfg.AccountStore.CreateSession(u.ID)
	if err != nil {
		t.Fatal(err)
	}

	mux := Routes(cfg)
	statusOf := func(rec *httptest.ResponseRecorder) string {
		var view phases.PhaseDetailView
		if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
			t.Fatal(err)
		}
		for _, stage := range view.Stages {
			for _, lab := range stage.Labs {
				if lab.ID == "alpha" {
					return lab.Status
				}
			}
		}
		t.Fatal("lab alpha not found in phase detail")
		return ""
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/phases/phase-1", nil))
	if rec.Code != http.StatusOK || statusOf(rec) != progress.LabStatusAvailable {
		t.Fatalf("phase detail with no session: code=%d status=%q, want %q",
			rec.Code, statusOf(rec), progress.LabStatusAvailable)
	}

	req := httptest.NewRequest("GET", "/api/phases/phase-1", nil)
	req.AddCookie(&http.Cookie{Name: webcore.SessionCookieName, Value: token})
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || statusOf(rec) != progress.LabStatusComplete {
		t.Fatalf("phase detail with session: code=%d status=%q, want %q (reflecting the synced attempt)",
			rec.Code, statusOf(rec), progress.LabStatusComplete)
	}
}
