package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/web/dashboard"
	"breaklab/internal/web/webcore"
)

func TestDashboardAPIMatchesTemplateHandlerData(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.AccountStore = newTestAccountStore(t)
	cfg.FrontendDir = t.TempDir()

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

	// No session: hosted attempts never count without one.
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/dashboard", nil))
	var noSession dashboard.DashboardView
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &noSession) != nil {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body)
	}
	if noSession.LabsCompleted != 0 {
		t.Errorf("no-session LabsCompleted = %d, want 0 (hosted attempt shouldn't count)", noSession.LabsCompleted)
	}

	// With session: the hosted attempt counts, same as the account-aware
	// template handler (TestDashboardAccountAwareReadPath).
	req := httptest.NewRequest("GET", "/api/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: webcore.SessionCookieName, Value: token})
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var withSession dashboard.DashboardView
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &withSession) != nil {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body)
	}
	if withSession.LabsCompleted != 1 {
		t.Errorf("with-session LabsCompleted = %d, want 1", withSession.LabsCompleted)
	}
}

func TestDashboardAPINotMountedWithoutFrontendDir(t *testing.T) {
	cfg := testConfig(t, 0) // FrontendDir left empty

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/dashboard", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404 (no FrontendDir configured)", rec.Code)
	}
}
