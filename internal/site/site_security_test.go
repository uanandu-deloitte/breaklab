package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A hand-edited progress.json that points a complete attempt's postmortem
// at a publish: true file outside postmortems/ must contribute nothing.
func TestBuildIgnoresPostmortemPathOutsideDir(t *testing.T) {
	cfg := newFixture(t)
	root := filepath.Dir(cfg.LabsDir)

	outside := filepath.Join(root, "elsewhere", "2026-10-17-web-down-config-typo.md")
	if err := os.MkdirAll(filepath.Dir(outside), 0o755); err != nil {
		t.Fatal(err)
	}
	body := strings.Replace(postmortemBody, "homepage stopped responding", "OUTSIDE_FILE_MARKER", 1)
	if err := os.WriteFile(outside, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfg.ProgressPath, []byte(sprintfProgress(outside)), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	var out strings.Builder
	_ = filepath.WalkDir(cfg.OutDir, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			data, _ := os.ReadFile(p)
			out.Write(data)
		}
		return nil
	})
	if strings.Contains(out.String(), "OUTSIDE_FILE_MARKER") {
		t.Error("content from a postmortem outside postmortems/ reached site/")
	}
	// The postmortems/ directory itself now always exists (it holds the
	// flat postmortems index page, which renders even when empty) — the
	// real invariant is that no page was written for the out-of-dir
	// postmortem specifically, and nothing links to it.
	if _, err := os.Stat(filepath.Join(cfg.OutDir, "postmortems", slug(outside)+".html")); err == nil {
		t.Error("a page was written for an out-of-dir postmortem path")
	}
	// index.html legitimately links to postmortems/index.html itself
	// (the "All postmortems" link); what must never appear is a link to
	// a per-postmortem page for the out-of-dir path specifically.
	if strings.Contains(out.String(), `href="postmortems/`+slug(outside)+`.html"`) || strings.Contains(out.String(), `href="`+slug(outside)+`.html"`) {
		t.Error("something links to a postmortem page for an out-of-dir path")
	}

	pmIndex := readFile(t, filepath.Join(cfg.OutDir, "postmortems", "index.html"))
	if len(strings.TrimSpace(pmIndex)) == 0 {
		t.Error("postmortems/index.html should still render (empty) when the only recorded postmortem is out-of-dir")
	}
}
