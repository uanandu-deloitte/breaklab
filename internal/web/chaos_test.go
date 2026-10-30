package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"breaklab/internal/web/chaos"
	"breaklab/internal/web/webtest"
)

func TestChaosAPIGroupsLabsByConcept(t *testing.T) {
	dir := t.TempDir()
	webtest.MustWriteFile(t, filepath.Join(dir, "roadmap.yaml"), "stages:\n  - id: linux\n    title: Linux fundamentals\n")
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "zombie-pileup", "lab.yaml"), `id: zombie-pileup
title: Zombie pileup
stage: linux
difficulty: 2
estimated_minutes: 20
tags: [fault-injection]
environment:
  dockerfile: Dockerfile
  ports: []
check:
  script: check.sh
  run_in: host
`)
	webtest.MustWriteFile(t, filepath.Join(dir, "progress", "progress.json"), `{"version": 1, "attempts": []}`)

	cfg := Config{
		LabsDir:      filepath.Join(dir, "labs"),
		RoadmapPath:  filepath.Join(dir, "roadmap.yaml"),
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		SessionPath:  filepath.Join(dir, ".breaklab", "session.json"),
		TemplatesDir: t.TempDir(),
		FrontendDir:  t.TempDir(),
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/chaos", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var view chaos.ChaosView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	var faultInjection chaos.ChaosConceptView
	for _, c := range view.Concepts {
		if c.ID == "fault-injection" {
			faultInjection = c
		}
	}
	if len(faultInjection.Labs) != 1 || faultInjection.Labs[0].ID != "zombie-pileup" {
		t.Fatalf("fault-injection labs = %+v, want [zombie-pileup]", faultInjection.Labs)
	}
}
