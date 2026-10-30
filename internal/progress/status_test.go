package progress

import (
	"testing"
	"time"

	"breaklab/internal/catalog"
)

func fixtureCatalogForStatus() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "complete-lab"},
			{ID: "active-lab"},
			{ID: "solved-lab"},
			{ID: "abandoned-lab"},
			{ID: "untouched-lab"},
			{ID: "abandoned-then-active-lab"},
		}},
	}}
}

func TestApplyStatus(t *testing.T) {
	now := time.Now()
	cat := fixtureCatalogForStatus()
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "complete-lab", Status: StatusComplete, SolvedAt: solvedPtr(now)},
		{Lab: "active-lab", Status: StatusActive},
		{Lab: "solved-lab", Status: StatusSolved, SolvedAt: solvedPtr(now)},
		{Lab: "abandoned-lab", Status: StatusAbandoned},
		{Lab: "abandoned-then-active-lab", Status: StatusAbandoned},
		{Lab: "abandoned-then-active-lab", Status: StatusActive},
	}}

	ApplyStatus(cat, pf)

	want := map[string]string{
		"complete-lab":              LabStatusComplete,
		"active-lab":                LabStatusInProgress,
		"solved-lab":                LabStatusInProgress,
		"abandoned-lab":             LabStatusAvailable,
		"untouched-lab":             LabStatusAvailable,
		"abandoned-then-active-lab": LabStatusInProgress,
	}
	for _, lab := range cat.Stages[0].Labs {
		if got, ok := want[lab.ID]; ok && lab.Status != got {
			t.Errorf("%s: status = %q, want %q", lab.ID, lab.Status, got)
		}
	}
}

func TestApplyStatusLocksLabWithUnmetPrerequisite(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "lnx-001"},
			{ID: "lnx-002", Prerequisites: []string{"lnx-001"}},
		}},
	}}

	// No attempts at all yet: lnx-001 is available, lnx-002 is locked
	// behind it.
	ApplyStatus(cat, &ProgressFile{})
	if got := statusOf(cat, "lnx-001"); got != LabStatusAvailable {
		t.Errorf("lnx-001 = %q, want available", got)
	}
	if got := statusOf(cat, "lnx-002"); got != LabStatusLocked {
		t.Errorf("lnx-002 = %q, want locked", got)
	}

	// Once lnx-001 is complete, lnx-002 unlocks to available.
	now := time.Now()
	cat = &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "lnx-001"},
			{ID: "lnx-002", Prerequisites: []string{"lnx-001"}},
		}},
	}}
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "lnx-001", Status: StatusComplete, SolvedAt: solvedPtr(now)},
	}}
	ApplyStatus(cat, pf)
	if got := statusOf(cat, "lnx-002"); got != LabStatusAvailable {
		t.Errorf("lnx-002 after lnx-001 complete = %q, want available", got)
	}
}

func TestApplyStatusNeverLocksAlreadyCompleteLab(t *testing.T) {
	// A lab that was already completed before its prerequisite existed
	// (or was ever re-validated) must stay complete, not regress to
	// locked.
	now := time.Now()
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "lnx-001"},
			{ID: "lnx-002", Prerequisites: []string{"lnx-001"}},
		}},
	}}
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "lnx-002", Status: StatusComplete, SolvedAt: solvedPtr(now)},
	}}
	ApplyStatus(cat, pf)
	if got := statusOf(cat, "lnx-002"); got != LabStatusComplete {
		t.Errorf("lnx-002 = %q, want complete (never locked once complete)", got)
	}
}

func statusOf(cat *catalog.Catalog, id string) string {
	for _, group := range cat.Stages {
		for _, lab := range group.Labs {
			if lab.ID == id {
				return lab.Status
			}
		}
	}
	return ""
}
