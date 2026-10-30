package postmortems

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

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

func TestBuildPostmortemPageRendersRegardlessOfPublishFlag(t *testing.T) {
	path := filepath.Join(t.TempDir(), "2026-10-17-web-down-config-typo.md")
	if err := os.WriteFile(path, []byte(postmortemFixtureBody), 0o644); err != nil {
		t.Fatal(err)
	}

	view, err := buildPostmortemPage(path, fixtureCatalogForPostmortemsIndex())
	if err != nil {
		t.Fatal(err)
	}
	if view.Lab != "The website returns nothing" || view.Variant != "config-typo" || view.Date != "2026-10-17" {
		t.Errorf("view = %+v", view)
	}
	if len(view.Sections) != 5 {
		t.Fatalf("sections = %+v", view.Sections)
	}
	for _, s := range view.Sections {
		if len(s.Blocks) == 0 {
			t.Errorf("section %q has no blocks", s.Name)
		}
	}
}

func TestBuildPostmortemPageShowsHypothesisWhenPresent(t *testing.T) {
	body := `---
lab: web-down
variant: config-typo
stage: linux
date: 2026-10-17
minutes_to_fix: 18
hints_used: 1
publish: false
---

# The website returns nothing

## Hypothesis
I bet the deploy script left the web server config in a half-applied state.

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
	path := filepath.Join(t.TempDir(), "2026-10-17-web-down-config-typo.md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	view, err := buildPostmortemPage(path, fixtureCatalogForPostmortemsIndex())
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Sections) != 6 {
		t.Fatalf("sections = %+v, want 6 (Hypothesis + the 5 required)", view.Sections)
	}
	if view.Sections[0].Name != "Hypothesis" {
		t.Fatalf("Sections[0].Name = %q, want Hypothesis first", view.Sections[0].Name)
	}
	if view.Sections[0].Blocks[0].Text != "I bet the deploy script left the web server config in a half-applied state." {
		t.Errorf("Hypothesis text = %q", view.Sections[0].Blocks[0].Text)
	}
}

func TestAttemptIndexForPostmortemPath(t *testing.T) {
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "a", Status: progress.StatusSolved},
		{Lab: "web-down", Status: progress.StatusComplete, Postmortem: "postmortems/2026-10-17-web-down-config-typo.md"},
	}}
	idx, ok := attemptIndexForPostmortemPath(pf, "postmortems/2026-10-17-web-down-config-typo.md")
	if !ok || idx != 1 {
		t.Errorf("attemptIndexForPostmortemPath = (%d, %v), want (1, true)", idx, ok)
	}
	if _, ok := attemptIndexForPostmortemPath(pf, "postmortems/no-such.md"); ok {
		t.Error("expected ok = false for an unknown path")
	}
}

func TestPostmortemSlugsOnlyIncludesCompleteAttempts(t *testing.T) {
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Status: progress.StatusComplete, Postmortem: "postmortems/2026-10-17-web-down-config-typo.md"},
		{Lab: "disk-full", Status: progress.StatusSolved, Postmortem: ""},
	}}

	slugs := postmortemSlugs(pf)
	if len(slugs) != 1 {
		t.Fatalf("slugs = %v, want exactly one", slugs)
	}
	if slugs["2026-10-17-web-down-config-typo"] != "postmortems/2026-10-17-web-down-config-typo.md" {
		t.Errorf("slugs = %v", slugs)
	}
}

func fixtureCatalogForPostmortemsIndex() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{
			Stage: catalog.Stage{ID: "linux", Title: "Linux fundamentals"},
			Labs: []catalog.Lab{
				{ID: "web-down", Title: "The website returns nothing", Difficulty: 1},
				{ID: "disk-full", Title: "Disk is full", Difficulty: 2},
			},
		},
		{
			Stage: catalog.Stage{ID: "git", Title: "Git"},
			Labs: []catalog.Lab{
				{ID: "git-lost-work", Title: "Lost commit", Difficulty: 1},
			},
		},
	}}
}

// writePostmortemFixture writes a minimal valid postmortem file with the
// given publish flag and returns its path, for tests that need
// isPublished to read real front matter from disk.
func writePostmortemFixture(t *testing.T, dir, name string, publish bool) string {
	t.Helper()
	path := filepath.Join(dir, name)
	body := strings.Replace(postmortemFixtureBody, "publish: false", "publish: "+strconv.FormatBool(publish), 1)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func fixtureProgressForPostmortemsIndex(t *testing.T, dir string) *progress.ProgressFile {
	t.Helper()
	webDownSolved := time.Date(2026, 10, 17, 10, 20, 0, 0, time.UTC)
	diskFullSolved := time.Date(2026, 10, 18, 9, 0, 0, 0, time.UTC)

	return &progress.ProgressFile{Attempts: []progress.Attempt{
		{
			Lab: "web-down", Variant: "config-typo",
			StartedAt: webDownSolved.Add(-18 * time.Minute), SolvedAt: &webDownSolved,
			Status:     progress.StatusComplete,
			Postmortem: writePostmortemFixture(t, dir, "web-down.md", true),
		},
		{
			Lab: "disk-full", Variant: "runaway-log",
			StartedAt: diskFullSolved.Add(-10 * time.Minute), SolvedAt: &diskFullSolved,
			Status:     progress.StatusComplete,
			Postmortem: writePostmortemFixture(t, dir, "disk-full.md", false),
		},
	}}
}

func TestBuildPostmortemsIndexListsAllRegardlessOfPublishStatus(t *testing.T) {
	dir := t.TempDir()
	cat := fixtureCatalogForPostmortemsIndex()
	pf := fixtureProgressForPostmortemsIndex(t, dir)

	view := buildPostmortemsIndex(cat, pf, url.Values{})
	if len(view.Rows) != 2 {
		t.Fatalf("Rows = %+v, want 2", view.Rows)
	}
	// Most recently solved first.
	if view.Rows[0].Lab != "Disk is full" || view.Rows[0].Status != "draft" {
		t.Errorf("Rows[0] = %+v, want Disk is full / draft", view.Rows[0])
	}
	if view.Rows[1].Lab != "The website returns nothing" || view.Rows[1].Status != "published" {
		t.Errorf("Rows[1] = %+v, want The website returns nothing / published", view.Rows[1])
	}
	if view.Rows[1].TicketID == "" {
		t.Error("expected a non-empty TicketID")
	}
}

func TestBuildPostmortemsIndexFiltersByStage(t *testing.T) {
	dir := t.TempDir()
	cat := fixtureCatalogForPostmortemsIndex()
	pf := fixtureProgressForPostmortemsIndex(t, dir)

	linux := buildPostmortemsIndex(cat, pf, url.Values{"stage": {"linux"}})
	if len(linux.Rows) != 2 {
		t.Errorf("stage=linux Rows = %+v, want 2", linux.Rows)
	}

	git := buildPostmortemsIndex(cat, pf, url.Values{"stage": {"git"}})
	if len(git.Rows) != 0 {
		t.Errorf("stage=git Rows = %+v, want 0 (no postmortem recorded for git-lost-work)", git.Rows)
	}
}

func TestBuildPostmortemsIndexFiltersByDifficulty(t *testing.T) {
	dir := t.TempDir()
	cat := fixtureCatalogForPostmortemsIndex()
	pf := fixtureProgressForPostmortemsIndex(t, dir)

	view := buildPostmortemsIndex(cat, pf, url.Values{"difficulty": {"2"}})
	if len(view.Rows) != 1 || view.Rows[0].Lab != "Disk is full" {
		t.Errorf("difficulty=2 Rows = %+v, want exactly Disk is full", view.Rows)
	}
}

func TestBuildPostmortemsIndexFiltersByStatus(t *testing.T) {
	dir := t.TempDir()
	cat := fixtureCatalogForPostmortemsIndex()
	pf := fixtureProgressForPostmortemsIndex(t, dir)

	published := buildPostmortemsIndex(cat, pf, url.Values{"status": {"published"}})
	if len(published.Rows) != 1 || published.Rows[0].Lab != "The website returns nothing" {
		t.Errorf("status=published Rows = %+v", published.Rows)
	}

	draft := buildPostmortemsIndex(cat, pf, url.Values{"status": {"draft"}})
	if len(draft.Rows) != 1 || draft.Rows[0].Lab != "Disk is full" {
		t.Errorf("status=draft Rows = %+v", draft.Rows)
	}
}

func TestBuildPostmortemsIndexFiltersCombinedAndClear(t *testing.T) {
	dir := t.TempDir()
	cat := fixtureCatalogForPostmortemsIndex()
	pf := fixtureProgressForPostmortemsIndex(t, dir)

	none := buildPostmortemsIndex(cat, pf, url.Values{"stage": {"linux"}, "status": {"published"}, "difficulty": {"2"}})
	if len(none.Rows) != 0 {
		t.Errorf("combined filters Rows = %+v, want 0 (Disk is full is draft, not published)", none.Rows)
	}

	cleared := buildPostmortemsIndex(cat, pf, url.Values{})
	if len(cleared.Rows) != 2 {
		t.Errorf("cleared filters Rows = %+v, want the full list of 2", cleared.Rows)
	}
}

func TestBuildPostmortemsIndexEmptyWhenNoPostmortems(t *testing.T) {
	cat := fixtureCatalogForPostmortemsIndex()
	view := buildPostmortemsIndex(cat, &progress.ProgressFile{}, url.Values{})
	if len(view.Rows) != 0 {
		t.Errorf("Rows = %+v, want empty", view.Rows)
	}
}
