package postmortem

import (
	"strings"
	"testing"
)

func TestSplitAndParseFrontMatter(t *testing.T) {
	content := "---\n" +
		"lab: web-down\n" +
		"variant: config-typo\n" +
		"stage: linux\n" +
		"date: 2026-10-17\n" +
		"minutes_to_fix: 18\n" +
		"hints_used: 1\n" +
		"publish: false\n" +
		"---\n" +
		"\n# The website returns nothing\n\n## Summary\nsome body text\n"

	raw, body, err := SplitFrontMatter(content)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "# The website returns nothing") || !strings.Contains(body, "## Summary") {
		t.Errorf("body = %q", body)
	}

	fm, err := ParseFrontMatter(raw)
	if err != nil {
		t.Fatal(err)
	}
	want := FrontMatter{
		Lab: "web-down", Variant: "config-typo", Stage: "linux",
		Date: "2026-10-17", MinutesToFix: 18, HintsUsed: 1, Publish: false,
	}
	if fm != want {
		t.Errorf("fm = %+v, want %+v", fm, want)
	}
}

func TestSplitFrontMatterMissingBlock(t *testing.T) {
	if _, _, err := SplitFrontMatter("# No front matter here\n"); err == nil {
		t.Fatal("expected an error for a file with no front matter")
	}
}

func TestSplitFrontMatterUnterminated(t *testing.T) {
	if _, _, err := SplitFrontMatter("---\nlab: web-down\n"); err == nil {
		t.Fatal("expected an error for a missing closing ---")
	}
}
