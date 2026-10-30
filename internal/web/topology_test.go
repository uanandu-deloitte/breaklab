package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"breaklab/internal/web/topology"
	"breaklab/internal/web/webtest"
)

// web-ui-design/web-ui spec, "Web UI attack surface covers Docker-backed
// routes": the container name passed to Docker must come only from the
// session file loaded server-side, never from the request.
func TestTopologyAPIIgnoresContainerFromRequest(t *testing.T) {
	cfg := topologyFixtureConfig(t, true)

	var gotName string
	cfg.DockerStatus = func(ctx context.Context, name string) (string, string, error) {
		gotName = name
		return "running", "breaklab-web-down", nil
	}

	req := httptest.NewRequest("GET", "/api/topology?container=breaklab-evil-other-lab", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if gotName != "breaklab-web-down" {
		t.Fatalf("Docker was queried for %q, want the session's own container breaklab-web-down regardless of the request's ?container= value", gotName)
	}
	if strings.Contains(rec.Body.String(), "breaklab-evil-other-lab") {
		t.Fatalf("response echoed the request-supplied container name")
	}
}

// web-ui spec's "same rule as every other view": no variant id ever
// reaches the response, matching the lab-detail rule.
func TestTopologyAPINeverLeaksVariantID(t *testing.T) {
	cfg := topologyFixtureConfig(t, true)
	cfg.DockerStatus = func(ctx context.Context, name string) (string, string, error) {
		return "running", "breaklab-web-down", nil
	}

	req := httptest.NewRequest("GET", "/api/topology", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "config-typo") {
		t.Fatalf("response leaked the active session's variant id")
	}
}

func TestTopologyAPINoActiveLabEmptyState(t *testing.T) {
	cfg := topologyFixtureConfig(t, false)

	req := httptest.NewRequest("GET", "/api/topology", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var view topology.TopologyView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.Active {
		t.Errorf("expected Active = false with no session, got %+v", view)
	}
}

// topologyFixtureConfig lays out a minimal on-disk fixture (roadmap + one
// lab) and, when active is true, an on-disk session for that lab.
// Topology migrated to the `frontend` capability's JSON API
// (GET /api/topology), so FrontendDir just needs to be non-empty for
// that route to be mounted — these tests never read its shell file.
func topologyFixtureConfig(t *testing.T, active bool) Config {
	t.Helper()
	dir := t.TempDir()

	webtest.MustWriteFile(t, filepath.Join(dir, "roadmap.yaml"), "stages:\n  - id: linux\n    title: Linux fundamentals\n")
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
variants:
  - id: config-typo
    script: variants/config-typo.sh
    weight: 1
`)
	webtest.MustWriteFile(t, filepath.Join(dir, "progress", "progress.json"), `{"version": 1, "attempts": []}`)

	sessionPath := filepath.Join(dir, ".breaklab", "session.json")
	if active {
		webtest.MustWriteFile(t, sessionPath, `{"lab":"web-down","container":"breaklab-web-down","variant":"config-typo","hints_used":0,"attempt_index":-1}`)
	}

	return Config{
		LabsDir:      filepath.Join(dir, "labs"),
		RoadmapPath:  filepath.Join(dir, "roadmap.yaml"),
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		SessionPath:  sessionPath,
		TemplatesDir: t.TempDir(),
		FrontendDir:  t.TempDir(),
	}
}
