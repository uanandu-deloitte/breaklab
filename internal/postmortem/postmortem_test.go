package postmortem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func copyTemplate(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "templates", "postmortem.md"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "postmortem.md")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRender(t *testing.T) {
	tmplPath := copyTemplate(t)
	fm := FrontMatter{
		Lab:          "web-down",
		Variant:      "config-typo",
		Stage:        "linux",
		Date:         "2026-10-17",
		MinutesToFix: 18,
		HintsUsed:    1,
		Publish:      false,
	}

	out, err := Render(tmplPath, fm, "The website returns nothing", "")
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"lab: web-down",
		"variant: config-typo",
		"stage: linux",
		"date: 2026-10-17",
		"minutes_to_fix: 18",
		"hints_used: 1",
		"publish: false",
		"# The website returns nothing",
		"## Hypothesis",
		"## Summary",
		"## Detection",
		"## Root cause",
		"## Fix",
		"## Prevention",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("rendered output missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "What did you think was broken") {
		t.Error("empty hypothesis should fall back to the placeholder comment")
	}
}

func TestRenderFillsInHypothesis(t *testing.T) {
	tmplPath := copyTemplate(t)
	fm := FrontMatter{Lab: "web-down", Variant: "config-typo", Stage: "linux", Date: "2026-10-17"}

	out, err := Render(tmplPath, fm, "The website returns nothing", "Probably the nginx config has a syntax error.")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Probably the nginx config has a syntax error.") {
		t.Errorf("rendered output missing the hypothesis text:\n%s", out)
	}
}
