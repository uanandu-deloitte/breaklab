package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"breaklab/internal/web/lab"
	"breaklab/internal/web/webtest"
)

func TestLabDetailAPIReturns404ForUnknownLab(t *testing.T) {
	dir := t.TempDir()
	webtest.MustWriteFile(t, filepath.Join(dir, "roadmap.yaml"), "stages:\n  - id: linux\n    title: Linux\n")
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "web-down", "lab.yaml"), `id: web-down
title: The website returns nothing
stage: linux
difficulty: 1
estimated_minutes: 20
environment:
  dockerfile: Dockerfile
  ports: ["8080:80"]
check:
  script: check.sh
  run_in: host
`)
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "web-down", "TICKET.md"), "ticket text\n")
	cfg := Config{
		LabsDir:      filepath.Join(dir, "labs"),
		RoadmapPath:  filepath.Join(dir, "roadmap.yaml"),
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		SessionPath:  filepath.Join(dir, ".breaklab", "session.json"),
		TemplatesDir: t.TempDir(),
		FrontendDir:  t.TempDir(),
	}

	req := httptest.NewRequest("GET", "/api/labs/no-such-lab", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

// TestLabsShellServedRegardlessOfLabExistence pins design.md's "Astro
// builds shells, not content" decision for this route: the HTML layer
// ("/labs/{id}") never checks whether id names a real lab — only the
// JSON API does (TestLabDetailAPIReturns404ForUnknownLab above).
func TestLabsShellServedRegardlessOfLabExistence(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()
	if err := os.MkdirAll(filepath.Join(cfg.FrontendDir, "labs", "_"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.FrontendDir, "labs", "_", "index.html"), []byte("<html>lab shell</html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, id := range []string{"alpha", "no-such-lab"} {
		rec := httptest.NewRecorder()
		Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/labs/"+id, nil))
		if rec.Code != http.StatusOK || rec.Body.String() != "<html>lab shell</html>" {
			t.Errorf("/labs/%s: code=%d body=%q", id, rec.Code, rec.Body.String())
		}
	}
}

// TestLabDetailAPIChallengeMode is TestLabDetailChallengeModeHidesHintsBlock's
// counterpart for the `frontend` capability: "/labs/{id}" itself is a
// static Astro shell now, so ChallengeMode (and the resulting
// hints-section visibility) lives in the JSON API response instead of
// server-rendered markup — the Astro page's own client-side code hides
// the Hints section/button when this field is true.
func TestLabDetailAPIChallengeMode(t *testing.T) {
	cfg := fullFixtureConfig(t)

	req := httptest.NewRequest("GET", "/api/labs/web-down?mode=challenge", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var view lab.LabDetailView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if !view.ChallengeMode {
		t.Error("?mode=challenge: ChallengeMode = false, want true")
	}

	req = httptest.NewRequest("GET", "/api/labs/web-down", nil)
	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	view = lab.LabDetailView{}
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.ChallengeMode {
		t.Error("no mode param: ChallengeMode = true, want false")
	}
}
