package progress

import (
	"testing"

	"breaklab/internal/catalog"
)

func patternFixtureCatalog() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "k8s"}, Labs: []catalog.Lab{
			{ID: "net-1", Tags: []string{"networking"}},
			{ID: "net-2", Tags: []string{"networking"}},
			{ID: "net-3", Tags: []string{"networking"}},
			{ID: "lnx-1", Tags: []string{"linux"}},
			{ID: "lnx-2", Tags: []string{"linux"}},
			{ID: "lnx-3", Tags: []string{"linux"}},
		}},
	}}
}

func complete(lab string, hints int) Attempt {
	return Attempt{Lab: lab, Status: StatusComplete, HintsUsed: hints}
}

func TestPatternsSingleQualifyingTag(t *testing.T) {
	cat := patternFixtureCatalog()
	pf := &ProgressFile{Attempts: []Attempt{
		complete("net-1", 6), complete("net-2", 6), complete("net-3", 6),
		complete("lnx-1", 1), complete("lnx-2", 1), complete("lnx-3", 1),
	}}

	got := Patterns(cat, pf)
	if len(got) != 1 || got[0].Tag != "networking" {
		t.Fatalf("got %+v, want exactly one pattern for networking", got)
	}
	if got[0].Count != 3 || got[0].AverageHintsUsed != 6 {
		t.Errorf("got %+v, want Count=3 AverageHintsUsed=6", got[0])
	}
}

func TestPatternsTagBelowSampleSizeFloorExcluded(t *testing.T) {
	cat := patternFixtureCatalog()
	pf := &ProgressFile{Attempts: []Attempt{
		complete("net-1", 10), complete("net-2", 10), // only 2 samples
		complete("lnx-1", 1), complete("lnx-2", 1), complete("lnx-3", 1),
	}}

	for _, p := range Patterns(cat, pf) {
		if p.Tag == "networking" {
			t.Fatalf("networking has only 2 samples, should be excluded: %+v", p)
		}
	}
}

func TestPatternsTagBelowMultiplierExcluded(t *testing.T) {
	cat := patternFixtureCatalog()
	// networking averages 2 hints, linux averages 1 — overall average is
	// 1.5, and 2 is below patternMultiplier (1.5) * 1.5 = 2.25.
	pf := &ProgressFile{Attempts: []Attempt{
		complete("net-1", 2), complete("net-2", 2), complete("net-3", 2),
		complete("lnx-1", 1), complete("lnx-2", 1), complete("lnx-3", 1),
	}}

	if got := Patterns(cat, pf); len(got) != 0 {
		t.Errorf("got %+v, want no patterns: 2 is below 1.5x the overall 1.5 average", got)
	}
}

func TestPatternsExactlyAtMultiplierThresholdQualifies(t *testing.T) {
	cat := patternFixtureCatalog()
	// networking averages 3, linux averages 1 — overall average is 2,
	// and 3 == patternMultiplier (1.5) * 2 exactly.
	pf := &ProgressFile{Attempts: []Attempt{
		complete("net-1", 3), complete("net-2", 3), complete("net-3", 3),
		complete("lnx-1", 1), complete("lnx-2", 1), complete("lnx-3", 1),
	}}

	got := Patterns(cat, pf)
	if len(got) != 1 || got[0].Tag != "networking" {
		t.Fatalf("got %+v, want networking to qualify exactly at the threshold", got)
	}
}

func TestPatternsExactlyAtSampleFloorQualifies(t *testing.T) {
	cat := patternFixtureCatalog()
	pf := &ProgressFile{Attempts: []Attempt{
		complete("net-1", 6), complete("net-2", 6), complete("net-3", 6), // exactly 3
		complete("lnx-1", 1), complete("lnx-2", 1), complete("lnx-3", 1),
	}}

	got := Patterns(cat, pf)
	if len(got) != 1 || got[0].Count != 3 {
		t.Fatalf("got %+v, want exactly 3 samples to qualify", got)
	}
}

func TestPatternsNoCompleteAttemptsReturnsNil(t *testing.T) {
	cat := patternFixtureCatalog()
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "net-1", Status: StatusActive},
	}}
	if got := Patterns(cat, pf); got != nil {
		t.Errorf("got %+v, want nil", got)
	}
}

func TestPatternsZeroOverallAverageReturnsNil(t *testing.T) {
	cat := patternFixtureCatalog()
	pf := &ProgressFile{Attempts: []Attempt{
		complete("net-1", 0), complete("net-2", 0), complete("net-3", 0),
	}}
	if got := Patterns(cat, pf); got != nil {
		t.Errorf("got %+v, want nil: nobody used any hints, nothing to be markedly above", got)
	}
}
