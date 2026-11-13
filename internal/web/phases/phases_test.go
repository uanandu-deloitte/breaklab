package phases

import (
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

func solvedAt(t time.Time) *time.Time { return &t }

func fixtureCatalogForPhases() *catalog.Catalog {
	linux := catalog.StageGroup{
		Stage: catalog.Stage{ID: "linux", Phase: "phase-1", Title: "Linux fundamentals"},
		Labs: []catalog.Lab{
			{ID: "complete-lab", Title: "Complete lab", Difficulty: 1, EstimatedMinutes: 20},
			{ID: "active-lab", Title: "Active lab", Difficulty: 2, EstimatedMinutes: 15},
		},
	}
	capstone := catalog.StageGroup{
		Stage: catalog.Stage{ID: "capstone", Phase: "phase-1", Title: "Capstone", Capstone: true},
		Labs:  []catalog.Lab{{ID: "final-incident", Title: "Everything is on fire"}},
	}
	return &catalog.Catalog{
		Stages: []catalog.StageGroup{linux, capstone},
		Phases: []catalog.PhaseGroup{
			{Phase: catalog.Phase{ID: "phase-1", Title: "Phase 1"}, Stages: []catalog.StageGroup{linux, capstone}},
			{Phase: catalog.Phase{ID: "phase-2", Title: "Phase 2 — Cloud", Blurb: "Coming soon."}, Stages: nil},
		},
	}
}

func TestBuildPhaseListAggregatesProgressAndIncludesEmptyPhase(t *testing.T) {
	cat := fixtureCatalogForPhases()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "complete-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}

	view := BuildPhaseList(cat, pf)
	if len(view.Phases) != 2 {
		t.Fatalf("len(Phases) = %d, want 2", len(view.Phases))
	}

	p1 := view.Phases[0]
	if p1.ID != "phase-1" || p1.LabsCompleted != 1 || p1.LabsTotal != 2 {
		t.Fatalf("phase-1 summary = %+v, want 1/2 complete (capstone excluded)", p1)
	}

	p2 := view.Phases[1]
	if p2.ID != "phase-2" || p2.LabsTotal != 0 || p2.LabsCompleted != 0 {
		t.Fatalf("phase-2 summary = %+v, want an empty, present phase", p2)
	}
}

func TestBuildPhaseDetailExcludesCapstoneAndSetsUpNext(t *testing.T) {
	cat := fixtureCatalogForPhases()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "complete-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}

	view, found := buildPhaseDetail(cat, pf, "phase-1")
	if !found {
		t.Fatal("expected phase-1 to be found")
	}
	if len(view.Stages) != 1 || view.Stages[0].ID != "linux" {
		t.Fatalf("expected only the regular (non-capstone) linux stage, got %+v", view.Stages)
	}

	var upNext []string
	for _, row := range view.Stages[0].Labs {
		if row.IsUpNext {
			upNext = append(upNext, row.ID)
		}
	}
	if len(upNext) != 1 || upNext[0] != "active-lab" {
		t.Fatalf("up-next rows = %v, want exactly [active-lab]", upNext)
	}
}

func TestBuildPhaseDetailEmptyPhaseRendersWithNoStages(t *testing.T) {
	cat := fixtureCatalogForPhases()
	view, found := buildPhaseDetail(cat, &progress.ProgressFile{}, "phase-2")
	if !found {
		t.Fatal("expected phase-2 to be found even with zero stages")
	}
	if view.Title != "Phase 2 — Cloud" || len(view.Stages) != 0 {
		t.Fatalf("expected phase-2's empty state, got %+v", view)
	}
}

func TestBuildPhaseDetailUnknownPhaseNotFound(t *testing.T) {
	cat := fixtureCatalogForPhases()
	_, found := buildPhaseDetail(cat, &progress.ProgressFile{}, "phase-nope")
	if found {
		t.Fatal("expected an unknown phase id to report not found")
	}
}
