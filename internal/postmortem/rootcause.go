package postmortem

import (
	"regexp"
	"strings"
)

// stopwords are common English words excluded from the root-cause
// consistency heuristic even though they clear the 4-character floor —
// without this, "with," "from," "were" etc. would count as "meaningful"
// overlap between any two unrelated sentences.
var stopwords = map[string]bool{
	"with": true, "from": true, "that": true, "this": true, "have": true,
	"been": true, "were": true, "will": true, "your": true, "when": true,
	"than": true, "then": true, "they": true, "what": true, "which": true,
	"while": true, "about": true, "into": true, "only": true, "does": true,
	"doesn": true, "isn": true, "after": true, "before": true, "there": true,
	"here": true, "some": true, "such": true, "most": true, "more": true,
}

var wordPattern = regexp.MustCompile(`[a-zA-Z]+`)

// MeaningfulWords lowercases text, extracts alphabetic tokens, and
// drops stopwords and anything under 4 characters — the small
// vocabulary RootCauseMentionsVariant compares against. Order is
// preserved but duplicates are removed.
func MeaningfulWords(text string) []string {
	var words []string
	seen := make(map[string]bool)
	for _, w := range wordPattern.FindAllString(strings.ToLower(text), -1) {
		if len(w) < 4 || stopwords[w] || seen[w] {
			continue
		}
		seen[w] = true
		words = append(words, w)
	}
	return words
}

// RootCauseMentionsVariant reports whether rootCause shares at least
// one meaningful word with hints' combined text — a cheap, explainable
// nudge that a postmortem's claimed root cause overlaps with what the
// solved variant's own hints describe as actually broken. It is a word
// -overlap heuristic, not a semantic correctness check: a true root
// cause described in entirely different words will not be detected,
// by design (see design.md - Decisions).
func RootCauseMentionsVariant(rootCause string, hints []string) bool {
	expected := make(map[string]bool)
	for _, w := range MeaningfulWords(strings.Join(hints, " ")) {
		expected[w] = true
	}
	for _, w := range MeaningfulWords(rootCause) {
		if expected[w] {
			return true
		}
	}
	return false
}
