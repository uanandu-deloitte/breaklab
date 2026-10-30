package challenges

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
	"breaklab/internal/web/webtest"
)

func solvedAt(t time.Time) *time.Time { return &t }

func fixtureCatalogForRoadmap() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{
			Stage: catalog.Stage{ID: "linux", Title: "Linux fundamentals"},
			Labs: []catalog.Lab{
				{ID: "complete-lab", Title: "Complete lab", Summary: "A one-line summary.", Difficulty: 1, EstimatedMinutes: 20, Tags: []string{"web-servers"}},
				{ID: "active-lab", Title: "Active lab", Difficulty: 2, EstimatedMinutes: 15, Tags: []string{"permissions"}},
				{ID: "untouched-lab", Title: "Untouched lab", Difficulty: 3, EstimatedMinutes: 30, Tags: []string{"web-servers"}},
			},
		},
		{
			Stage: catalog.Stage{ID: "git", Title: "Git", Blurb: "recover from disaster"},
			Labs: []catalog.Lab{
				{ID: "git-lab", Title: "Git lab", Difficulty: 1, EstimatedMinutes: 10, Tags: []string{"networking"}},
			},
		},
	}}
}

func TestBuildChallengesNoFiltersReturnsFullList(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{}, nil)
	if len(view.Rows) != 4 {
		t.Fatalf("got %d rows, want 4 (every lab across every stage)", len(view.Rows))
	}
}

func TestBuildChallengesFiltersByStage(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"stage": {"git"}}, nil)
	if len(view.Rows) != 1 || view.Rows[0].ID != "git-lab" {
		t.Fatalf("got %+v, want only git-lab", view.Rows)
	}
}

func TestBuildChallengesFiltersByDifficulty(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"difficulty": {"2"}}, nil)
	for _, row := range view.Rows {
		if row.Difficulty != 2 {
			t.Errorf("row %+v has difficulty != 2", row)
		}
	}
	if len(view.Rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(view.Rows))
	}
}

func TestBuildChallengesFiltersByStatus(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "complete-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"status": {progress.LabStatusComplete}}, nil)
	if len(view.Rows) != 1 || view.Rows[0].ID != "complete-lab" {
		t.Fatalf("got %+v, want only complete-lab", view.Rows)
	}
}

func TestBuildChallengesUnrecognizedFilterYieldsEmptyNotError(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	cases := []url.Values{
		{"stage": {"no-such-stage"}},
		{"status": {"no-such-status"}},
		{"difficulty": {"not-a-number"}},
	}
	for _, q := range cases {
		view := buildChallenges(context.Background(), cat, pf, nil, q, nil)
		if len(view.Rows) != 0 {
			t.Errorf("query %v: got %d rows, want 0", q, len(view.Rows))
		}
	}
}

func TestBuildChallengesStagesListedInRoadmapOrder(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	view := buildChallenges(context.Background(), cat, &progress.ProgressFile{}, nil, url.Values{}, nil)
	if len(view.Stages) != 2 || view.Stages[0].ID != "linux" || view.Stages[1].ID != "git" {
		t.Errorf("Stages = %+v", view.Stages)
	}
}

func TestBuildChallengesFiltersByTag(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"tag": {"networking"}}, nil)
	if len(view.Rows) != 1 || view.Rows[0].ID != "git-lab" {
		t.Fatalf("got %+v, want only git-lab", view.Rows)
	}
}

func TestBuildChallengesTagsListedSortedAndDeduped(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	view := buildChallenges(context.Background(), cat, &progress.ProgressFile{}, nil, url.Values{}, nil)
	want := []string{"networking", "permissions", "web-servers"}
	if len(view.Tags) != len(want) {
		t.Fatalf("Tags = %v, want %v", view.Tags, want)
	}
	for i, tag := range want {
		if view.Tags[i] != tag {
			t.Fatalf("Tags = %v, want %v", view.Tags, want)
		}
	}
}

func TestBuildChallengesSearchesTitleAndTags(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"q": {"COMPLETE"}}, nil)
	if len(view.Rows) != 1 || view.Rows[0].ID != "complete-lab" {
		t.Fatalf("title search: got %+v, want only complete-lab", view.Rows)
	}

	view = buildChallenges(context.Background(), cat, pf, nil, url.Values{"q": {"networking"}}, nil)
	if len(view.Rows) != 1 || view.Rows[0].ID != "git-lab" {
		t.Fatalf("tag search: got %+v, want only git-lab", view.Rows)
	}

	view = buildChallenges(context.Background(), cat, pf, nil, url.Values{"q": {"no-such-text"}}, nil)
	if len(view.Rows) != 0 {
		t.Fatalf("got %d rows, want 0 for a non-matching search", len(view.Rows))
	}
}

func TestBuildChallengesSearchCombinesWithFilters(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	// "web-servers" tag matches both complete-lab and untouched-lab;
	// narrowing by stage should still apply alongside the search.
	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"q": {"web-servers"}, "stage": {"linux"}}, nil)
	if len(view.Rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(view.Rows))
	}
}

func TestBuildChallengesSortDifficulty(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"sort": {"difficulty"}}, nil)
	for i := 1; i < len(view.Rows); i++ {
		if view.Rows[i-1].Difficulty > view.Rows[i].Difficulty {
			t.Fatalf("Rows not sorted by difficulty: %+v", view.Rows)
		}
	}
}

func TestBuildChallengesSortAlpha(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{"sort": {"alpha"}}, nil)
	for i := 1; i < len(view.Rows); i++ {
		if view.Rows[i-1].Title > view.Rows[i].Title {
			t.Fatalf("Rows not sorted alphabetically: %+v", view.Rows)
		}
	}
}

func TestBuildChallengesSortRecentDefault(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "complete-lab", Status: progress.StatusComplete, StartedAt: base, SolvedAt: solvedAt(base.Add(time.Hour))},
		{Lab: "active-lab", Status: progress.StatusActive, StartedAt: base.Add(24 * time.Hour)},
	}}

	view := buildChallenges(context.Background(), cat, pf, nil, url.Values{}, nil)
	// active-lab started more recently than complete-lab, and both more
	// recently than labs with no attempt at all (zero time, sorts last).
	var activeIdx, completeIdx, untouchedIdx = -1, -1, -1
	for i, row := range view.Rows {
		switch row.ID {
		case "active-lab":
			activeIdx = i
		case "complete-lab":
			completeIdx = i
		case "untouched-lab":
			untouchedIdx = i
		}
	}
	if !(activeIdx < completeIdx && completeIdx < untouchedIdx) {
		t.Fatalf("recent-activity order wrong: active=%d complete=%d untouched=%d", activeIdx, completeIdx, untouchedIdx)
	}
}

func TestBuildChallengesActiveIncidentBanner(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "active-lab", Status: progress.StatusActive, StartedAt: time.Now()},
	}}
	session := &runner.Session{Lab: "active-lab", AttemptIndex: 0}

	view := buildChallenges(context.Background(), cat, pf, session, url.Values{}, nil)
	if view.ActiveIncident == nil {
		t.Fatal("ActiveIncident = nil, want a banner")
	}
	if view.ActiveIncident.Lab != "active-lab" || view.ActiveIncident.TicketID != webcore.TicketID(0) {
		t.Errorf("ActiveIncident = %+v", view.ActiveIncident)
	}
}

func TestBuildChallengesNoActiveIncident(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	view := buildChallenges(context.Background(), cat, &progress.ProgressFile{}, nil, url.Values{}, nil)
	if view.ActiveIncident != nil {
		t.Errorf("ActiveIncident = %+v, want nil", view.ActiveIncident)
	}
}

func TestChallengesRandomRedirectsToAMatchingLab(t *testing.T) {
	dir := t.TempDir()
	webtest.MustWriteFile(t, filepath.Join(dir, "roadmap.yaml"), "stages:\n  - id: git\n    title: Git\n")
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "git-lab", "lab.yaml"), `id: git-lab
title: Git lab
stage: git
difficulty: 1
estimated_minutes: 10
environment:
  dockerfile: Dockerfile
check:
  script: check.sh
  run_in: host
`)
	templatesDir := t.TempDir()

	cfg := webcore.Config{
		LabsDir:      filepath.Join(dir, "labs"),
		RoadmapPath:  filepath.Join(dir, "roadmap.yaml"),
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		TemplatesDir: templatesDir,
	}

	req := httptest.NewRequest("POST", "/api/challenges/random", nil)
	rec := httptest.NewRecorder()
	challengesRandomAPIHandler(cfg)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var payload randomChallengePayload
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Lab != "git-lab" {
		t.Errorf("lab = %q, want git-lab", payload.Lab)
	}
}

func TestChallengesRandomHonorsFormBodyFilters(t *testing.T) {
	dir := t.TempDir()
	webtest.MustWriteFile(t, filepath.Join(dir, "roadmap.yaml"), "stages:\n  - id: git\n    title: Git\n  - id: linux\n    title: Linux\n")
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "git-lab", "lab.yaml"), `id: git-lab
title: Git lab
stage: git
difficulty: 1
estimated_minutes: 10
environment:
  dockerfile: Dockerfile
check:
  script: check.sh
  run_in: host
`)
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "linux-lab", "lab.yaml"), `id: linux-lab
title: Linux lab
stage: linux
difficulty: 1
estimated_minutes: 10
environment:
  dockerfile: Dockerfile
check:
  script: check.sh
  run_in: host
`)
	templatesDir := t.TempDir()
	cfg := webcore.Config{
		LabsDir:      filepath.Join(dir, "labs"),
		RoadmapPath:  filepath.Join(dir, "roadmap.yaml"),
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		TemplatesDir: templatesDir,
	}

	req := httptest.NewRequest("POST", "/api/challenges/random?stage=linux", nil)
	rec := httptest.NewRecorder()
	challengesRandomAPIHandler(cfg)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var payload randomChallengePayload
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Lab != "linux-lab" {
		t.Errorf("lab = %q, want linux-lab (stage=linux filter)", payload.Lab)
	}
}

func TestChallengesRandomAPINoMatchReturnsEmptyLab(t *testing.T) {
	dir := t.TempDir()
	webtest.MustWriteFile(t, filepath.Join(dir, "roadmap.yaml"), "stages:\n  - id: git\n    title: Git\n")
	webtest.MustWriteFile(t, filepath.Join(dir, "labs", "git-lab", "lab.yaml"), `id: git-lab
title: Git lab
stage: git
difficulty: 1
estimated_minutes: 10
environment:
  dockerfile: Dockerfile
check:
  script: check.sh
  run_in: host
`)
	templatesDir := t.TempDir()

	cfg := webcore.Config{
		LabsDir:      filepath.Join(dir, "labs"),
		RoadmapPath:  filepath.Join(dir, "roadmap.yaml"),
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		TemplatesDir: templatesDir,
	}

	req := httptest.NewRequest("POST", "/api/challenges/random?stage=no-such-stage", nil)
	rec := httptest.NewRecorder()
	challengesRandomAPIHandler(cfg)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var payload randomChallengePayload
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Lab != "" {
		t.Errorf("lab = %q, want empty (no match)", payload.Lab)
	}
}
