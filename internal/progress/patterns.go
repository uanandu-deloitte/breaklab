package progress

import (
	"sort"

	"breaklab/internal/catalog"
)

// patternMinSamples and patternMultiplier are the fixed thresholds a
// tag must clear to be reported as a pattern — see design.md's
// Decisions for why these aren't user-configurable in v1.
const (
	patternMinSamples = 3
	patternMultiplier = 1.5
)

// Pattern is one topic tag whose average hints_used, across its own
// complete attempts, is markedly above the user's overall average.
type Pattern struct {
	Tag              string
	AverageHintsUsed float64
	Count            int
}

// Patterns returns every tag meeting the pattern threshold (at least
// patternMinSamples complete attempts, average hints_used at least
// patternMultiplier times the overall average across every complete
// attempt), sorted by AverageHintsUsed descending. Returns nil when
// there's no complete attempt yet, or when the overall average is zero
// (nothing to be "markedly above").
func Patterns(cat *catalog.Catalog, pf *ProgressFile) []Pattern {
	tagsByLab := make(map[string][]string)
	for _, group := range cat.RegularStages() {
		for _, lab := range group.Labs {
			tagsByLab[lab.ID] = lab.Tags
		}
	}

	type sums struct{ hints, count int }
	perTag := make(map[string]*sums)
	var overallHints, overallCount int

	for _, a := range pf.Attempts {
		if a.Status != StatusComplete {
			continue
		}
		overallHints += a.HintsUsed
		overallCount++
		for _, tag := range tagsByLab[a.Lab] {
			s, ok := perTag[tag]
			if !ok {
				s = &sums{}
				perTag[tag] = s
			}
			s.hints += a.HintsUsed
			s.count++
		}
	}

	if overallCount == 0 || overallHints == 0 {
		return nil
	}
	overallAvg := float64(overallHints) / float64(overallCount)

	var patterns []Pattern
	for tag, s := range perTag {
		if s.count < patternMinSamples {
			continue
		}
		avg := float64(s.hints) / float64(s.count)
		if avg >= patternMultiplier*overallAvg {
			patterns = append(patterns, Pattern{Tag: tag, AverageHintsUsed: avg, Count: s.count})
		}
	}
	sort.Slice(patterns, func(i, j int) bool {
		if patterns[i].AverageHintsUsed != patterns[j].AverageHintsUsed {
			return patterns[i].AverageHintsUsed > patterns[j].AverageHintsUsed
		}
		return patterns[i].Tag < patterns[j].Tag
	})
	return patterns
}
