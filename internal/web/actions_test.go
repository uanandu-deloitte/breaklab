package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/runner"
)

// fakeActionDocker is a minimal runner.Docker fake scoped to this
// package's tests; internal/runner has its own, but it's unexported in
// a different package.
type fakeActionDocker struct {
	checkOut []byte
	checkErr error
}

func (f *fakeActionDocker) Ping(context.Context) error { return nil }
func (f *fakeActionDocker) Build(context.Context, string, string, string) error {
	return nil
}
func (f *fakeActionDocker) Run(context.Context, string, string, []string, bool, []string, int) error {
	return nil
}
func (f *fakeActionDocker) Exec(context.Context, string, ...string) ([]byte, error) {
	return f.checkOut, f.checkErr
}
func (f *fakeActionDocker) Copy(context.Context, string, string, string) error { return nil }
func (f *fakeActionDocker) Remove(context.Context, string) error               { return nil }
func (f *fakeActionDocker) VolumeRemove(context.Context, string) error         { return nil }
func (f *fakeActionDocker) VolumeCreate(context.Context, string) error         { return nil }

const actionRoadmapYAML = "stages:\n  - id: linux\n    title: Linux fundamentals\n"

func actionLabYAML(id string) string {
	return "id: " + id + `
title: Test lab
stage: linux
difficulty: 1
estimated_minutes: 5
environment:
  dockerfile: Dockerfile
  ports: []
check:
  script: check.sh
  run_in: container
hints:
  - general hint one
variants:
  - id: only-variant
    script: variants/only.sh
    weight: 1
    hints:
      - variant hint
`
}

// actionTestConfig builds a Config whose every path lives under a fresh
// t.TempDir(), proving action routes never touch the real repo's
// labs/, progress/, or .breaklab/ even though the actions package builds
// a fresh *runner.Runner per request.
func actionTestConfig(t *testing.T, fd *fakeActionDocker) Config {
	t.Helper()
	root := t.TempDir()
	labDir := filepath.Join(root, "labs", "alpha")
	if err := os.MkdirAll(filepath.Join(labDir, "variants"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"lab.yaml":         actionLabYAML("alpha"),
		"Dockerfile":       "FROM scratch\n",
		"TICKET.md":        "# INC: alpha is down\n",
		"check.sh":         "#!/usr/bin/env bash\necho PASS: fine\n",
		"variants/only.sh": "#!/usr/bin/env bash\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(labDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "roadmap.yaml"), []byte(actionRoadmapYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	templatesDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(templatesDir, "lab.html"), []byte("{{.Title}}"), 0o644); err != nil {
		t.Fatal(err)
	}
	frontendDir := t.TempDir()

	var docker runner.Docker
	if fd != nil {
		docker = fd
	}
	return Config{
		LabsDir:      filepath.Join(root, "labs"),
		RoadmapPath:  filepath.Join(root, "roadmap.yaml"),
		SessionPath:  filepath.Join(root, ".breaklab", "session.json"),
		HistoryPath:  filepath.Join(root, ".breaklab", "history.json"),
		ProgressPath: filepath.Join(root, "progress", "progress.json"),
		TemplatesDir: templatesDir,
		FrontendDir:  frontendDir,
		Docker:       docker,
	}
}

// writeActionLab adds a second lab fixture (id) under cfg.LabsDir,
// identical in shape to the "alpha" lab actionTestConfig already wrote,
// so a test can exercise lab-switch or multi-account scenarios.
func writeActionLab(t *testing.T, cfg Config, id string) {
	t.Helper()
	dir := filepath.Join(cfg.LabsDir, id)
	if err := os.MkdirAll(filepath.Join(dir, "variants"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"lab.yaml":         actionLabYAML(id),
		"Dockerfile":       "FROM scratch\n",
		"TICKET.md":        "# INC: " + id + " is down\n",
		"check.sh":         "#!/usr/bin/env bash\necho PASS: fine\n",
		"variants/only.sh": "#!/usr/bin/env bash\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// hostedActionTestConfig is actionTestConfig with a fresh AccountStore
// attached, so action routes require a logged-in session instead of
// running unauthenticated.
func hostedActionTestConfig(t *testing.T, fd *fakeActionDocker) Config {
	t.Helper()
	cfg := actionTestConfig(t, fd)
	cfg.AccountStore = newTestAccountStore(t)
	return cfg
}

// apiErrorPayload mirrors webcore's unexported type of the same name —
// the one JSON error shape every API response uses — so these
// black-box tests can decode it without needing webcore to export it.
type apiErrorPayload struct {
	Error string `json:"error"`
}

// hintPayload mirrors actions' unexported type of the same name.
type hintPayload struct {
	Hint string `json:"hint"`
}

// decodeJSONError reads an apiErrorPayload from rec's body, failing the
// test if it doesn't parse.
func decodeJSONError(t *testing.T, rec *httptest.ResponseRecorder) apiErrorPayload {
	t.Helper()
	var payload apiErrorPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decoding error body %q: %v", rec.Body.String(), err)
	}
	return payload
}

func TestActionAPIRoutesNilDockerReturnUnavailable(t *testing.T) {
	cfg := actionTestConfig(t, nil) // no Docker configured

	for _, path := range []string{"/api/labs/alpha/start", "/api/labs/alpha/check", "/api/labs/alpha/hint", "/api/labs/alpha/hypothesis", "/api/labs/alpha/stop"} {
		req := httptest.NewRequest("POST", path, nil)
		rec := httptest.NewRecorder()
		Routes(cfg).ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("%s: status = %d, want %d (body: %s)", path, rec.Code, http.StatusServiceUnavailable, rec.Body.String())
			continue
		}
		if got := decodeJSONError(t, rec).Error; !strings.Contains(got, "not available") {
			t.Errorf("%s: error = %q, want it to mention unavailability", path, got)
		}
	}
}

func TestLabLifecycleThroughActionAPIRoutes(t *testing.T) {
	fd := &fakeActionDocker{checkOut: []byte("PASS: storefront healthy\n")}
	cfg := actionTestConfig(t, fd)
	t.Setenv("BREAKLAB_FORCE_VARIANT", "only-variant")

	// Start: no body needed, session now points at this lab.
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/start", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("start status = %d, want 204 (body: %s)", rec.Code, rec.Body.String())
	}
	session, err := runner.LoadSession(cfg.SessionPath)
	if err != nil || session == nil || session.Lab != "alpha" {
		t.Fatalf("session after start = %+v, %v", session, err)
	}

	// Hint: reveals the lab's one general hint.
	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/hint", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("hint status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var hint hintPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &hint); err != nil {
		t.Fatalf("decoding hint body %q: %v", rec.Body.String(), err)
	}
	if hint.Hint != "general hint one" {
		t.Errorf("hint = %q, want %q", hint.Hint, "general hint one")
	}
	session, _ = runner.LoadSession(cfg.SessionPath)
	if session.HintsUsed != 1 {
		t.Errorf("HintsUsed after one hint = %d, want 1", session.HintsUsed)
	}

	// Check: fake Docker reports PASS, the response carries the message.
	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/check", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("check status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var res runner.CheckResult
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decoding check body %q: %v", rec.Body.String(), err)
	}
	if res.Message != "PASS: storefront healthy" {
		t.Errorf("check message = %q", res.Message)
	}
	if !res.Passed {
		t.Error("Passed missing for a passing check")
	}

	// Stop: session is cleared.
	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/stop", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("stop status = %d, want 204 (body: %s)", rec.Code, rec.Body.String())
	}
	session, err = runner.LoadSession(cfg.SessionPath)
	if err != nil || session != nil {
		t.Fatalf("session after stop = %+v, %v, want nil", session, err)
	}
}

func TestLabHintRefusedDuringReviewAttempt(t *testing.T) {
	fd := &fakeActionDocker{checkOut: []byte("PASS: fine\n")}
	cfg := actionTestConfig(t, fd)
	t.Setenv("BREAKLAB_FORCE_VARIANT", "only-variant")

	// Seed a prior complete attempt for "alpha" before starting a new one.
	solved := time.Now().Add(-24 * time.Hour)
	pf := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "alpha", Variant: "only-variant", SolvedAt: &solved, Status: progress.StatusComplete},
	}}
	if err := progress.Save(cfg.ProgressPath, pf); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/start", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("start status = %d, want 204 (body: %s)", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/hint", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("hint status = %d, want 500 (body: %s)", rec.Code, rec.Body.String())
	}
	if got := decodeJSONError(t, rec).Error; !strings.Contains(got, "review attempt") {
		t.Errorf("error = %q, want it to mention the review attempt", got)
	}
}

func TestLabHypothesisSavedViaAPI(t *testing.T) {
	fd := &fakeActionDocker{checkOut: []byte("PASS: fine\n")}
	cfg := actionTestConfig(t, fd)
	t.Setenv("BREAKLAB_FORCE_VARIANT", "only-variant")

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/start", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("start status = %d, want 204 (body: %s)", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/labs/alpha/hypothesis", strings.NewReader(`{"hypothesis":"probably a config typo"}`))
	req.Header.Set("Content-Type", "application/json")
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("hypothesis status = %d, want 204 (body: %s)", rec.Code, rec.Body.String())
	}

	session, _ := runner.LoadSession(cfg.SessionPath)
	pf, err := progress.Load(cfg.ProgressPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := pf.Attempts[session.AttemptIndex].Hypothesis; got != "probably a config typo" {
		t.Errorf("saved hypothesis = %q", got)
	}
}

func TestLabStartBlockedByOtherLabRequiresConfirmAPI(t *testing.T) {
	fd := &fakeActionDocker{checkOut: []byte("PASS: fine\n")}
	cfg := actionTestConfig(t, fd)
	// A second lab under the same root so Start has something real to switch to.
	writeActionLab(t, cfg, "beta")
	t.Setenv("BREAKLAB_FORCE_VARIANT", "only-variant")

	// Start alpha first.
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/alpha/start", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("start alpha status = %d, want 204 (body: %s)", rec.Code, rec.Body.String())
	}

	// Plain start of beta is refused, naming the blocker.
	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("POST", "/api/labs/beta/start", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("plain start beta status = %d, want 500 (body: %s)", rec.Code, rec.Body.String())
	}
	if got := decodeJSONError(t, rec).Error; !strings.Contains(got, "alpha") {
		t.Errorf("error = %q, want it to name the blocking lab alpha", got)
	}
	session, _ := runner.LoadSession(cfg.SessionPath)
	if session.Lab != "alpha" {
		t.Errorf("session after refused start = %q, want it still alpha", session.Lab)
	}

	// Confirmed start of beta switches over.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/labs/beta/start", strings.NewReader(`{"confirm":true}`))
	req.Header.Set("Content-Type", "application/json")
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("confirmed start beta status = %d, want 204 (body: %s)", rec.Code, rec.Body.String())
	}
	session, _ = runner.LoadSession(cfg.SessionPath)
	if session.Lab != "beta" {
		t.Errorf("session after confirmed start = %q, want beta", session.Lab)
	}
}

func TestLabActionsRequireLoginWhenHosted(t *testing.T) {
	fd := &fakeActionDocker{checkOut: []byte("PASS: fine\n")}
	cfg := hostedActionTestConfig(t, fd)
	mux := Routes(cfg)

	for _, path := range []string{"/api/labs/alpha/start", "/api/labs/alpha/check", "/api/labs/alpha/hint", "/api/labs/alpha/hypothesis", "/api/labs/alpha/stop"} {
		req := httptest.NewRequest("POST", path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s with no session: status = %d, want 401 (body: %s)", path, rec.Code, rec.Body)
			continue
		}
		if got := decodeJSONError(t, rec).Error; !strings.Contains(got, "not logged in") {
			t.Errorf("%s: error = %q, want it to mention not being logged in", path, got)
		}
	}
}
