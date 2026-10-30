package site

import (
	"os"
	"path/filepath"
	"testing"

	"breaklab/internal/progress"
)

func writePostmortem(t *testing.T, path string, publish bool) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	content := `---
lab: web-down
variant: config-typo
stage: linux
date: 2026-10-17
minutes_to_fix: 18
hints_used: 1
publish: ` + boolStr(publish) + `
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
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func TestLoadPublishedIncludesHypothesisWhenPresent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "with-hypothesis.md")
	content := `---
lab: web-down
variant: config-typo
stage: linux
date: 2026-10-17
minutes_to_fix: 18
hints_used: 1
publish: true
---

# The website returns nothing

## Hypothesis
I bet the deploy script left the config half-applied.

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
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, views, err := LoadPublished(dir, []string{path})
	if err != nil {
		t.Fatal(err)
	}
	view := views[path]
	if len(view.Sections) != 6 || view.Sections[0].Name != "Hypothesis" {
		t.Fatalf("sections = %+v, want Hypothesis first followed by the 5 required", view.Sections)
	}
	if view.Sections[0].Paragraphs[0] != "I bet the deploy script left the config half-applied." {
		t.Errorf("hypothesis paragraph = %q", view.Sections[0].Paragraphs[0])
	}
}

func TestCollectPostmortemPaths(t *testing.T) {
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Status: progress.StatusComplete, Postmortem: "postmortems/a.md"},
		{Lab: "disk-full", Status: progress.StatusSolved, Postmortem: ""},
		{Lab: "git-lost-work", Status: progress.StatusComplete, Postmortem: ""},
	}}
	got := CollectPostmortemPaths(pf)
	if len(got) != 1 || got[0] != "postmortems/a.md" {
		t.Errorf("got %v", got)
	}
}

func TestLoadPublishedIncludesPublishedOnly(t *testing.T) {
	dir := t.TempDir()
	publishedPath := filepath.Join(dir, "published.md")
	unpublishedPath := filepath.Join(dir, "unpublished.md")
	writePostmortem(t, publishedPath, true)
	writePostmortem(t, unpublishedPath, false)

	published, views, err := LoadPublished(dir, []string{publishedPath, unpublishedPath})
	if err != nil {
		t.Fatal(err)
	}

	if !published[publishedPath] {
		t.Error("published postmortem missing from published map")
	}
	if published[unpublishedPath] {
		t.Error("unpublished postmortem present in published map")
	}
	if _, ok := views[unpublishedPath]; ok {
		t.Error("unpublished postmortem has a view")
	}

	view, ok := views[publishedPath]
	if !ok {
		t.Fatal("published postmortem has no view")
	}
	if view.Lab != "web-down" || view.Variant != "config-typo" || view.Date != "2026-10-17" {
		t.Errorf("view = %+v", view)
	}
	if len(view.Sections) != 5 {
		t.Fatalf("sections = %+v", view.Sections)
	}
	for _, s := range view.Sections {
		if len(s.Paragraphs) == 0 {
			t.Errorf("section %q has no paragraphs", s.Name)
		}
	}
}

func TestLoadPublishedSkipsMissingFile(t *testing.T) {
	missingDir := t.TempDir()
	published, views, err := LoadPublished(missingDir, []string{filepath.Join(missingDir, "does-not-exist.md")})
	if err != nil {
		t.Fatal(err)
	}
	if len(published) != 0 || len(views) != 0 {
		t.Errorf("published = %v, views = %v, want both empty", published, views)
	}
}
