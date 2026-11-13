// Package postmortem renders and validates the postmortem a learner
// writes for a solved lab: templates/postmortem.md filled in with the
// attempt's facts, and the five required sections checked for real
// content before the attempt is promoted to "complete".
package postmortem

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// FrontMatter is a postmortem file's YAML front matter.
type FrontMatter struct {
	Lab          string `yaml:"lab"`
	Variant      string `yaml:"variant"`
	Stage        string `yaml:"stage"`
	Date         string `yaml:"date"`
	MinutesToFix int    `yaml:"minutes_to_fix"`
	HintsUsed    int    `yaml:"hints_used"`
	Publish      bool   `yaml:"publish"`
}

// renderData is FrontMatter plus the two fields the template needs that
// aren't part of the front matter itself: the lab's title, used for the
// file's # heading, and the learner's own pre-fix guess (if any), used
// to pre-fill the Hypothesis section.
type renderData struct {
	FrontMatter
	Title      string
	Hypothesis string
}

// Render fills in the template at templatePath with fm, title, and
// hypothesis (the learner's own guess at the cause, written while the
// attempt was still active — "" if they never recorded one).
func Render(templatePath string, fm FrontMatter, title, hypothesis string) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, renderData{FrontMatter: fm, Title: title, Hypothesis: hypothesis}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderSections builds a complete postmortem file from fm, the lab's
// title, and real content already written for each of DisplaySections
// (keyed by section name; a missing key renders as an empty section).
// Unlike Render, this writes the learner's actual words straight in —
// no HTML-comment placeholder prompts — for a caller (the web
// postmortem editor) that collects finished content directly rather
// than handing the learner a file to fill in by hand. fm's string
// fields are catalog-validated ids or values this package already
// trusts (lab/variant/stage ids, an RFC3339 date), so they're written
// as unquoted YAML scalars the same way templates/postmortem.md does.
func RenderSections(fm FrontMatter, title string, sections map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "---\nlab: %s\nvariant: %s\nstage: %s\ndate: %s\nminutes_to_fix: %d\nhints_used: %d\npublish: %t\n---\n\n# %s\n",
		fm.Lab, fm.Variant, fm.Stage, fm.Date, fm.MinutesToFix, fm.HintsUsed, fm.Publish, title)
	for _, name := range DisplaySections {
		fmt.Fprintf(&b, "\n## %s\n%s\n", name, strings.TrimSpace(sections[name]))
	}
	return b.String()
}
