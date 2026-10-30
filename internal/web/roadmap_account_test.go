package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/web/roadmap"
	"breaklab/internal/web/webcore"
)

func TestRoadmapAPIAccountAwareReadPath(t *testing.T) {
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
		var view roadmap.RoadmapView
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
		t.Fatal("lab alpha not found in roadmap")
		return ""
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/roadmap", nil))
	if rec.Code != http.StatusOK || statusOf(rec) != progress.LabStatusAvailable {
		t.Fatalf("roadmap with no session: code=%d status=%q, want %q (synced attempt from another device shouldn't apply)",
			rec.Code, statusOf(rec), progress.LabStatusAvailable)
	}

	req := httptest.NewRequest("GET", "/api/roadmap", nil)
	req.AddCookie(&http.Cookie{Name: webcore.SessionCookieName, Value: token})
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || statusOf(rec) != progress.LabStatusComplete {
		t.Fatalf("roadmap with session: code=%d status=%q, want %q (reflecting the synced attempt)",
			rec.Code, statusOf(rec), progress.LabStatusComplete)
	}
}
