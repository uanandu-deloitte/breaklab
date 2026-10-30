package site

import (
	"os"
	"strings"

	"breaklab/internal/postmortem"
	"breaklab/internal/progress"
)

// SectionView is one rendered section of a published postmortem.
type SectionView struct {
	Name       string
	Paragraphs []string
}

// PostmortemView is everything site/postmortems/<slug>.html renders.
type PostmortemView struct {
	Lab      string
	Variant  string
	Date     string
	Slug     string
	Sections []SectionView
}

// CollectPostmortemPaths returns the Postmortem field of every "complete"
// attempt in pf.
func CollectPostmortemPaths(pf *progress.ProgressFile) []string {
	var paths []string
	for _, a := range pf.Attempts {
		if a.Status == progress.StatusComplete && a.Postmortem != "" {
			paths = append(paths, a.Postmortem)
		}
	}
	return paths
}

// LoadPublished reads each path and, for the ones whose front matter
// says publish: true, records them in published and renders them into
// views. A publish: false, missing, unreadable, or out-of-dir file is skipped —
// it's absent from both maps, never causing the whole load to fail
// (design.md: "nothing leaks" is enforced here, not by erroring out).
func LoadPublished(dir string, paths []string) (published map[string]bool, views map[string]PostmortemView, err error) {
	published = make(map[string]bool)
	views = make(map[string]PostmortemView)

	for _, path := range paths {
		// A path outside dir (or a symlink) is skipped like a missing file.
		confined, confineErr := postmortem.Confined(dir, path)
		if confineErr != nil {
			continue
		}
		data, readErr := os.ReadFile(confined)
		if readErr != nil {
			continue
		}
		raw, body, splitErr := postmortem.SplitFrontMatter(string(data))
		if splitErr != nil {
			continue
		}
		fm, parseErr := postmortem.ParseFrontMatter(raw)
		if parseErr != nil || !fm.Publish {
			continue
		}

		published[path] = true
		views[path] = PostmortemView{
			Lab:      fm.Lab,
			Variant:  fm.Variant,
			Date:     fm.Date,
			Slug:     slug(path),
			Sections: buildSections(body),
		}
	}
	return published, views, nil
}

// buildSections renders each of postmortem.DisplaySections' text as
// paragraphs, in order, after stripping placeholder comments. Hypothesis
// is optional and skipped entirely when blank, unlike the required
// sections it's listed alongside.
func buildSections(body string) []SectionView {
	raw := postmortem.SplitIntoSections(body)
	sections := make([]SectionView, 0, len(postmortem.DisplaySections))
	for _, name := range postmortem.DisplaySections {
		text := postmortem.StripComments(raw[name])
		if name == "Hypothesis" && strings.TrimSpace(text) == "" {
			continue
		}
		sections = append(sections, SectionView{Name: name, Paragraphs: paragraphs(text)})
	}
	return sections
}

// paragraphs splits text on blank lines into non-empty, trimmed blocks.
func paragraphs(text string) []string {
	var out []string
	for _, block := range strings.Split(text, "\n\n") {
		if block = strings.TrimSpace(block); block != "" {
			out = append(out, block)
		}
	}
	return out
}
