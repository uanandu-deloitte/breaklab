package postmortem

import (
	"regexp"
	"strings"
)

// RequiredSections lists the five sections every postmortem must have,
// in the order templates/postmortem.md presents them.
var RequiredSections = []string{"Summary", "Detection", "Root cause", "Fix", "Prevention"}

// DisplaySections lists every section a postmortem viewer should render,
// in order: Hypothesis (the learner's pre-fix guess, written before the
// root cause was known) ahead of RequiredSections. Hypothesis is
// deliberately not part of RequiredSections or ValidateSections — it's
// optional, so a postmortem written before this field existed, or one
// where the learner skipped it, still validates exactly as before.
var DisplaySections = append([]string{"Hypothesis"}, RequiredSections...)

var htmlComment = regexp.MustCompile(`(?s)<!--.*?-->`)

// StripComments removes any HTML comments (including multi-line ones)
// from text, such as templates/postmortem.md's placeholder prompts.
func StripComments(text string) string {
	return htmlComment.ReplaceAllString(text, "")
}

// ValidateSections checks each of RequiredSections has at least 20 words
// of real content once placeholder HTML comments are stripped out, and
// returns the names of any that don't (an empty slice means valid). A
// missing heading counts as 0 words, so it fails the same way an
// untouched placeholder does.
func ValidateSections(body string) []string {
	sections := SplitIntoSections(body)

	var failing []string
	for _, name := range RequiredSections {
		text := StripComments(sections[name])
		if len(strings.Fields(text)) < 20 {
			failing = append(failing, name)
		}
	}
	return failing
}

// SplitIntoSections maps each "## <name>" heading in body to the text
// between it and the next "##" heading (or EOF).
func SplitIntoSections(body string) map[string]string {
	sections := make(map[string]string)
	lines := strings.Split(body, "\n")

	var current string
	var buf []string
	flush := func() {
		if current != "" {
			sections[current] = strings.Join(buf, "\n")
		}
	}
	for _, line := range lines {
		if name, ok := sectionHeading(line); ok {
			flush()
			current, buf = name, nil
			continue
		}
		if current != "" {
			buf = append(buf, line)
		}
	}
	flush()
	return sections
}

func sectionHeading(line string) (string, bool) {
	const prefix = "## "
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, prefix)), true
}
