package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/site"
	"breaklab/internal/web/webtest"
)

func publicProfileTestConfig(t *testing.T) Config {
	t.Helper()
	cfg := actionTestConfig(t, nil)
	cfg.AccountStore = newTestAccountStore(t)
	return cfg
}

// publicProfileView mirrors the account package's unexported type of
// the same name.
type publicProfileView struct {
	Username string
	Stages   []site.StageView
	Totals   site.Totals
}

func TestPublicProfileAPINotOptedInIs404(t *testing.T) {
	cfg := publicProfileTestConfig(t)
	if _, err := cfg.AccountStore.CreateUser("ada@example.com", "ada", "correct-password"); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/u/ada", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404 (account hasn't opted in)", rec.Code)
	}
}

func TestPublicProfileAPIUnknownUsernameIs404(t *testing.T) {
	cfg := publicProfileTestConfig(t)

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/u/nobody", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404", rec.Code)
	}
}

func TestPublicProfileAPIOptedInShowsProgress(t *testing.T) {
	cfg := publicProfileTestConfig(t)
	u, err := cfg.AccountStore.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.AccountStore.SetPublicProfile(u.ID, true); err != nil {
		t.Fatal(err)
	}
	solved := time.Now().UTC()
	if err := cfg.AccountStore.PushAttempts(u.ID, []progress.Attempt{
		{ID: "h1", Lab: "alpha", Status: progress.StatusComplete, StartedAt: solved.Add(-time.Minute), SolvedAt: &solved},
	}); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/u/ada", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body.String())
	}
	var view publicProfileView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.Username != "ada" || view.Totals.LabsCompleted != 1 {
		t.Fatalf("view = %+v, want ada with 1 completed lab", view)
	}
}

func TestPublicProfileAPIOptOutIs404Immediately(t *testing.T) {
	cfg := publicProfileTestConfig(t)
	u, err := cfg.AccountStore.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.AccountStore.SetPublicProfile(u.ID, true); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/u/ada", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("while opted in: code = %d, want 200", rec.Code)
	}

	if err := cfg.AccountStore.SetPublicProfile(u.ID, false); err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/u/ada", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("immediately after opting out: code = %d, want 404", rec.Code)
	}
}

// TestPublicProfileShellServedRegardlessOfOptIn: the /u/{name} page
// itself is a static Astro shell (the "Astro builds shells, not
// content" pattern) — it always 200s for any username; only the JSON
// API above encodes the opt-in/not-found distinction.
func TestPublicProfileShellServedRegardlessOfOptIn(t *testing.T) {
	cfg := publicProfileTestConfig(t)
	webtest.MustWriteFile(t, filepath.Join(cfg.FrontendDir, "u", "_", "index.html"), "<html></html>")

	for _, name := range []string{"ada", "nobody-at-all"} {
		rec := httptest.NewRecorder()
		Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/u/"+name, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("/u/%s: code = %d, want 200 (shell is served regardless of whether that account exists)", name, rec.Code)
		}
	}
}
