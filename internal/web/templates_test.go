package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/web/challenges"
	"breaklab/internal/web/dashboard"
	"breaklab/internal/web/roadmap"
	"breaklab/internal/web/webtest"
)

const realTemplatesDir = "../../templates/web"

// postmortemFixtureBody is a minimal valid postmortem file, shared by
// this file's fullFixtureConfig and security_test.go.
const postmortemFixtureBody = `---
lab: web-down
variant: config-typo
stage: linux
date: 2026-10-17
minutes_to_fix: 18
hints_used: 1
publish: false
---

# The website returns nothing

## Summary
The homepage stopped responding after a routine maintenance window and
customers saw connection errors for about twenty minutes this morning.

## Detection
I ran curl against the site and got a connection refused error, then
checked whether nginx was running at all with systemctl status nginx.

## Root cause
A trailing semicolon was missing from the listen directive in the nginx
site config, so nginx refused to start after the maintenance script ran.

## Fix
I added the missing semicolon back to the listen directive, ran nginx
dash t to check the config, and then started nginx again successfully.

## Prevention
Add a CI step that runs nginx config test on every change to the nginx
configuration files before they are allowed to merge into the main branch.
`

// fullFixtureConfig lays out a complete on-disk fixture (roadmap, a lab,
// a progress.json with one complete, published-false attempt and its
// postmortem file) and points TemplatesDir at the real templates/web/
// directory, so these tests exercise the actual shipped templates.
func fullFixtureConfig(t *testing.T) Config {
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
hints:
  - a general hint
variants:
  - id: config-typo
    script: variants/config-typo.sh
    weight: 1
`)
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "web-down", "TICKET.md"), "# INC-101\nWebsite is down.\n")

	postmortemPath := filepath.Join(dir, "postmortems", "2026-10-17-web-down-config-typo.md")
	webtest.MustWriteFile(t, postmortemPath, postmortemFixtureBody)

	webtest.MustWriteFile(t, filepath.Join(dir, "progress", "progress.json"), `{
  "version": 1,
  "attempts": [
    {
      "lab": "web-down",
      "variant": "config-typo",
      "started_at": "2026-10-17T10:02:00Z",
      "solved_at": "2026-10-17T10:20:00Z",
      "failed_checks": 1,
      "hints_used": 1,
      "postmortem": "`+postmortemPath+`",
      "status": "complete"
    }
  ]
}`)

	templatesDir, err := filepath.Abs(realTemplatesDir)
	if err != nil {
		t.Fatal(err)
	}

	// "/" and "/roadmap" have no html/template fallback anymore (the
	// `frontend` capability's Astro pages) — a minimal, script-free
	// static shell is enough for these tests, which care about every
	// *other* page's real template output, not those pages'.
	frontendDir := filepath.Join(dir, "frontend-dist")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "dashboard", "index.html"), "<html><body>dashboard</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "roadmap", "index.html"), "<html><body>roadmap</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "challenges", "index.html"), "<html><body>challenges</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "labs", "_", "index.html"), "<html><body>lab</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "topology", "index.html"), "<html><body>topology</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "postmortems", "index.html"), "<html><body>postmortems</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "postmortems", "_", "index.html"), "<html><body>postmortem</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "chaos", "index.html"), "<html><body>chaos</body></html>")
	webtest.MustWriteFile(t, filepath.Join(frontendDir, "capstone", "index.html"), "<html><body>capstone</body></html>")

	return Config{
		LabsDir:        filepath.Join(dir, "labs"),
		RoadmapPath:    filepath.Join(dir, "roadmap.yaml"),
		ProgressPath:   filepath.Join(dir, "progress", "progress.json"),
		SessionPath:    filepath.Join(dir, ".breaklab", "session.json"),
		TemplatesDir:   templatesDir,
		FrontendDir:    frontendDir,
		PostmortemsDir: filepath.Join(dir, "postmortems"),
	}
}

// TestRealDashboardAPIReportsPatterns exercises GET /api/dashboard
// (the dashboard's data now, since "/" itself is a static Astro
// shell with no embedded data): the base fixture has no tag with
// enough complete attempts to qualify as a pattern, so Patterns must
// be empty; adding three same-tag complete attempts past the
// threshold must make one appear.
func TestRealDashboardAPIReportsPatterns(t *testing.T) {
	cfg := fullFixtureConfig(t)

	req := httptest.NewRequest("GET", "/api/dashboard", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var view dashboard.DashboardView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if len(view.Patterns) != 0 {
		t.Errorf("expected no patterns: no tag has enough complete attempts to qualify, got %+v", view.Patterns)
	}

	for _, id := range []string{"net-1", "net-2", "net-3"} {
		webtest.MustWriteFile(t, filepath.Join(cfg.LabsDir, id, "lab.yaml"), tagLabYAML(id, "networking"))
	}
	for _, id := range []string{"lnx-1", "lnx-2", "lnx-3"} {
		webtest.MustWriteFile(t, filepath.Join(cfg.LabsDir, id, "lab.yaml"), tagLabYAML(id, "linux"))
	}
	webtest.MustWriteFile(t, cfg.ProgressPath, `{
  "version": 1,
  "attempts": [
    {"lab": "net-1", "variant": "v", "started_at": "2026-10-17T10:00:00Z", "solved_at": "2026-10-17T10:30:00Z", "hints_used": 6, "status": "complete"},
    {"lab": "net-2", "variant": "v", "started_at": "2026-10-17T10:00:00Z", "solved_at": "2026-10-17T10:30:00Z", "hints_used": 6, "status": "complete"},
    {"lab": "net-3", "variant": "v", "started_at": "2026-10-17T10:00:00Z", "solved_at": "2026-10-17T10:30:00Z", "hints_used": 6, "status": "complete"},
    {"lab": "lnx-1", "variant": "v", "started_at": "2026-10-17T10:00:00Z", "solved_at": "2026-10-17T10:30:00Z", "hints_used": 1, "status": "complete"},
    {"lab": "lnx-2", "variant": "v", "started_at": "2026-10-17T10:00:00Z", "solved_at": "2026-10-17T10:30:00Z", "hints_used": 1, "status": "complete"},
    {"lab": "lnx-3", "variant": "v", "started_at": "2026-10-17T10:00:00Z", "solved_at": "2026-10-17T10:30:00Z", "hints_used": 1, "status": "complete"}
  ]
}`)

	req = httptest.NewRequest("GET", "/api/dashboard", nil)
	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	view = dashboard.DashboardView{}
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range view.Patterns {
		if p.Tag == "networking" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a pattern naming networking, got %+v", view.Patterns)
	}
}

func tagLabYAML(id, tag string) string {
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
  run_in: host
tags: [` + tag + `]
`
}

// TestRealChallengesAPIReportsReviewDue is the `frontend`-capability
// counterpart of the old TestRealChallengesRendersReviewDueBadge:
// "/challenges" is a static Astro shell now, so this checks the JSON
// API's ReviewDue field instead of rendered badge markup.
func TestRealChallengesAPIReportsReviewDue(t *testing.T) {
	cfg := fullFixtureConfig(t)
	old := time.Now().Add(-(progress.ReviewInterval + time.Hour)).UTC().Format(time.RFC3339)
	webtest.MustWriteFile(t, cfg.ProgressPath, fmt.Sprintf(`{
  "version": 1,
  "attempts": [
    {
      "lab": "web-down",
      "variant": "config-typo",
      "started_at": %q,
      "solved_at": %q,
      "hints_used": 0,
      "status": "complete"
    }
  ]
}`, old, old))

	req := httptest.NewRequest("GET", "/api/challenges", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var view challenges.ChallengeView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, row := range view.Rows {
		if row.ID == "web-down" && row.ReviewDue {
			found = true
		}
	}
	if !found {
		t.Errorf("expected web-down to be ReviewDue, got %+v", view.Rows)
	}
}

// TestRealRoadmapAPIReportsReviewDue is TestRealChallengesRendersReviewDueBadge's
// counterpart for the roadmap, which moved to the `frontend` capability's
// JSON API ("/roadmap" itself is a static Astro shell with no embedded
// data now).
func TestRealRoadmapAPIReportsReviewDue(t *testing.T) {
	cfg := fullFixtureConfig(t)
	old := time.Now().Add(-(progress.ReviewInterval + time.Hour)).UTC().Format(time.RFC3339)
	webtest.MustWriteFile(t, cfg.ProgressPath, fmt.Sprintf(`{
  "version": 1,
  "attempts": [
    {
      "lab": "web-down",
      "variant": "config-typo",
      "started_at": %q,
      "solved_at": %q,
      "hints_used": 0,
      "status": "complete"
    }
  ]
}`, old, old))

	req := httptest.NewRequest("GET", "/api/roadmap", nil)
	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var view roadmap.RoadmapView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, stage := range view.Stages {
		for _, lab := range stage.Labs {
			if lab.ID == "web-down" && lab.ReviewDue {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expected web-down to be ReviewDue, got %+v", view.Stages)
	}
}
