package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const roadmapYAML = "stages:\n  - id: linux\n    title: Linux fundamentals\n"

const labYAML = `id: web-down
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
  - SECRET_HINT_TEXT must never reach site/
variants:
  - id: config-typo
    script: variants/config-typo.sh
    weight: 1
    hints:
      - SECRET_VARIANT_HINT_TEXT must never reach site/
`

const postmortemBody = `---
lab: web-down
variant: config-typo
stage: linux
date: 2026-10-17
minutes_to_fix: 18
hints_used: 1
publish: true
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

const progressJSON = `{
  "version": 1,
  "attempts": [
    {
      "lab": "web-down",
      "variant": "config-typo",
      "started_at": "2026-10-17T10:02:00Z",
      "solved_at": "2026-10-17T10:20:00Z",
      "failed_checks": 1,
      "hints_used": 1,
      "postmortem": "%s",
      "status": "complete"
    }
  ]
}
`

// newFixture lays out a temp repo with a lab (including variant/solution
// scripts and hints carrying secret markers that must never leak into
// site/), roadmap.yaml, progress.json, and a published postmortem.
func newFixture(t *testing.T) Config {
	t.Helper()
	dir := t.TempDir()

	write := func(rel, content string) {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("roadmap.yaml", roadmapYAML)
	write("labs/web-down/lab.yaml", labYAML)
	write("labs/web-down/Dockerfile", "FROM scratch\n")
	write("labs/web-down/TICKET.md", "# INC: web-down is down\n")
	write("labs/web-down/check.sh", "#!/usr/bin/env bash\necho PASS\n")
	write("labs/web-down/variants/config-typo.sh", "#!/usr/bin/env bash\n# SECRET_VARIANT_SCRIPT_CONTENT must never reach site/\n")
	write("labs/web-down/solutions/config-typo.sh", "#!/usr/bin/env bash\n# SECRET_SOLUTION_SCRIPT_CONTENT must never reach site/\n")

	postmortemPath := filepath.Join(dir, "postmortems", "2026-10-17-web-down-config-typo.md")
	write("postmortems/2026-10-17-web-down-config-typo.md", postmortemBody)
	write("progress/progress.json", sprintfProgress(postmortemPath))

	templatesDir, err := filepath.Abs(filepath.Join("..", "..", "templates", "site"))
	if err != nil {
		t.Fatal(err)
	}

	return Config{
		LabsDir:        filepath.Join(dir, "labs"),
		RoadmapPath:    filepath.Join(dir, "roadmap.yaml"),
		ProgressPath:   filepath.Join(dir, "progress", "progress.json"),
		TemplatesDir:   templatesDir,
		OutDir:         filepath.Join(dir, "site"),
		PostmortemsDir: filepath.Join(dir, "postmortems"),
	}
}

func sprintfProgress(postmortemPath string) string {
	escaped := strings.ReplaceAll(postmortemPath, `\`, `\\`)
	return strings.Replace(progressJSON, "%s", escaped, 1)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestBuildCreatesIndexAndPostmortemPage(t *testing.T) {
	cfg := newFixture(t)
	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	index := readFile(t, filepath.Join(cfg.OutDir, "index.html"))
	if !strings.Contains(index, "Linux fundamentals") || !strings.Contains(index, "The website returns nothing") {
		t.Errorf("index.html missing expected content:\n%s", index)
	}
	if !strings.Contains(index, `href="postmortems/2026-10-17-web-down-config-typo.html"`) {
		t.Errorf("index.html missing link to published postmortem:\n%s", index)
	}

	pmPath := filepath.Join(cfg.OutDir, "postmortems", "2026-10-17-web-down-config-typo.html")
	pm := readFile(t, pmPath)
	for _, want := range []string{"Summary", "Detection", "Root cause", "Fix", "Prevention", "homepage stopped responding"} {
		if !strings.Contains(pm, want) {
			t.Errorf("postmortem page missing %q:\n%s", want, pm)
		}
	}

	pmIndex := readFile(t, filepath.Join(cfg.OutDir, "postmortems", "index.html"))
	if !strings.Contains(pmIndex, "The website returns nothing") || !strings.Contains(pmIndex, `href="2026-10-17-web-down-config-typo.html"`) {
		t.Errorf("postmortems/index.html missing expected entry:\n%s", pmIndex)
	}

	if _, err := os.Stat(filepath.Join(cfg.OutDir, "style.css")); err != nil {
		t.Errorf("style.css not written: %v", err)
	}
	for _, name := range []string{"JetBrainsMono-Regular.woff2", "JetBrainsMono-Medium.woff2"} {
		if _, err := os.Stat(filepath.Join(cfg.OutDir, "fonts", name)); err != nil {
			t.Errorf("fonts/%s not written: %v", name, err)
		}
	}
}

func TestBuildUnpublishedProducesNoPage(t *testing.T) {
	cfg := newFixture(t)
	// Flip the postmortem to unpublished before the first (and only) build.
	pmPath := filepath.Join(filepath.Dir(cfg.ProgressPath), "..", "postmortems", "2026-10-17-web-down-config-typo.md")
	data := readFile(t, pmPath)
	data = strings.Replace(data, "publish: true", "publish: false", 1)
	if err := os.WriteFile(pmPath, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	index := readFile(t, filepath.Join(cfg.OutDir, "index.html"))
	if strings.Contains(index, "postmortems/2026-10-17-web-down-config-typo.html") {
		t.Errorf("index.html links to an unpublished postmortem:\n%s", index)
	}
	if _, err := os.Stat(filepath.Join(cfg.OutDir, "postmortems", "2026-10-17-web-down-config-typo.html")); err == nil {
		t.Error("unpublished postmortem page was written")
	}

	pmIndex := readFile(t, filepath.Join(cfg.OutDir, "postmortems", "index.html"))
	if strings.Contains(pmIndex, "The website returns nothing") {
		t.Errorf("postmortems/index.html lists an unpublished postmortem:\n%s", pmIndex)
	}
}

func TestBuildNothingLeaks(t *testing.T) {
	cfg := newFixture(t)
	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	var allOutput strings.Builder
	var sawVariantsOrSolutionsPath bool
	if err := filepath.WalkDir(cfg.OutDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.Contains(path, string(filepath.Separator)+"variants"+string(filepath.Separator)) ||
			strings.Contains(path, string(filepath.Separator)+"solutions"+string(filepath.Separator)) {
			sawVariantsOrSolutionsPath = true
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		allOutput.Write(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if sawVariantsOrSolutionsPath {
		t.Error("a variants/ or solutions/ path was copied into the built site")
	}
	for _, secret := range []string{
		"SECRET_HINT_TEXT", "SECRET_VARIANT_HINT_TEXT",
		"SECRET_VARIANT_SCRIPT_CONTENT", "SECRET_SOLUTION_SCRIPT_CONTENT",
	} {
		if strings.Contains(allOutput.String(), secret) {
			t.Errorf("%s leaked into the built site", secret)
		}
	}
}

func TestBuildWorksOffline(t *testing.T) {
	cfg := newFixture(t)
	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	index := readFile(t, filepath.Join(cfg.OutDir, "index.html"))
	pm := readFile(t, filepath.Join(cfg.OutDir, "postmortems", "2026-10-17-web-down-config-typo.html"))
	pmIndex := readFile(t, filepath.Join(cfg.OutDir, "postmortems", "index.html"))

	for name, html := range map[string]string{"index.html": index, "postmortem page": pm, "postmortems/index.html": pmIndex} {
		if strings.Contains(html, "http://") || strings.Contains(html, "https://") {
			t.Errorf("%s references an external URL", name)
		}
		if strings.Contains(html, "<script") {
			t.Errorf("%s includes a <script> tag", name)
		}
	}

	if !strings.Contains(index, `href="style.css"`) {
		t.Errorf("index.html's stylesheet link isn't the expected relative path:\n%s", index)
	}
	if !strings.Contains(pm, `href="../style.css"`) {
		t.Errorf("postmortem page's stylesheet link isn't the expected relative path:\n%s", pm)
	}
	if !strings.Contains(pmIndex, `href="../style.css"`) {
		t.Errorf("postmortems/index.html's stylesheet link isn't the expected relative path:\n%s", pmIndex)
	}
}

func TestBuildRemovesStalePageAfterUnpublishing(t *testing.T) {
	cfg := newFixture(t)
	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}
	pmOutPath := filepath.Join(cfg.OutDir, "postmortems", "2026-10-17-web-down-config-typo.html")
	if _, err := os.Stat(pmOutPath); err != nil {
		t.Fatalf("expected first build to publish the page: %v", err)
	}

	pmSourcePath := filepath.Join(filepath.Dir(cfg.ProgressPath), "..", "postmortems", "2026-10-17-web-down-config-typo.md")
	data := readFile(t, pmSourcePath)
	data = strings.Replace(data, "publish: true", "publish: false", 1)
	if err := os.WriteFile(pmSourcePath, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(pmOutPath); err == nil {
		t.Error("stale postmortem page survived a rebuild after unpublishing")
	}
}
