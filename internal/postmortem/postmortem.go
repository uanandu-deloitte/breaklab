// Package postmortem renders and validates the postmortem a learner
// writes for a solved lab: templates/postmortem.md filled in with the
// attempt's facts, and the five required sections checked for real
// content before the attempt is promoted to "complete".
package postmortem

import (
	"bytes"
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
