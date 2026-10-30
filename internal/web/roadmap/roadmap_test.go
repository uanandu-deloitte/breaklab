package roadmap

import (
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/web/webcore"
)

func solvedAt(t time.Time) *time.Time { return &t }

func fixtureCatalogForRoadmap() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{
			Stage: catalog.Stage{ID: "linux", Title: "Linux fundamentals"},
			Labs: []catalog.Lab{
				{ID: "complete-lab", Title: "Complete lab", Summary: "A one-line summary.", Difficulty: 1, EstimatedMinutes: 20, Tags: []string{"web-servers"}},
				{ID: "active-lab", Title: "Active lab", Difficulty: 2, EstimatedMinutes: 15, Tags: []string{"permissions"}},
				{ID: "untouched-lab", Title: "Untouched lab", Difficulty: 3, EstimatedMinutes: 30, Tags: []string{"web-servers"}},
			},
		},
		{
			Stage: catalog.Stage{ID: "git", Title: "Git", Blurb: "recover from disaster"},
			Labs: []catalog.Lab{
				{ID: "git-lab", Title: "Git lab", Difficulty: 1, EstimatedMinutes: 10, Tags: []string{"networking"}},
			},
		},
	}}
}

func TestBuildRoadmapStageChain(t *testing.T) {
	view := buildRoadmap(fixtureCatalogForRoadmap(), &progress.ProgressFile{})

	if len(view.Stages) != 2 {
		t.Fatalf("len(Stages) = %d, want 2", len(view.Stages))
	}
	linux, git := view.Stages[0], view.Stages[1]

	if linux.Index != 1 || git.Index != 2 {
		t.Errorf("Index = %d, %d, want 1, 2", linux.Index, git.Index)
	}
	if linux.IsLast {
		t.Error("first stage should not be IsLast")
	}
	if !git.IsLast {
		t.Error("last stage should be IsLast")
	}
	if linux.ConnectorBlurb != "recover from disaster" {
		t.Errorf("linux.ConnectorBlurb = %q, want the *next* stage's (git's) blurb", linux.ConnectorBlurb)
	}
	if git.ConnectorBlurb != "" {
		t.Errorf("git.ConnectorBlurb = %q, want empty: it's the last stage, no connector follows it", git.ConnectorBlurb)
	}
}

func TestBuildRoadmapAllThreeStatuses(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "complete-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now()), Postmortem: "postmortems/complete-lab.md"},
		{Lab: "active-lab", Status: progress.StatusActive},
	}}

	view := buildRoadmap(cat, pf)
	if len(view.Stages) != 2 || view.Stages[0].ID != "linux" || view.Stages[1].ID != "git" {
		t.Fatalf("stages out of roadmap order: %+v", view.Stages)
	}

	linux := view.Stages[0]
	byID := map[string]webcore.LabRowView{}
	for _, row := range linux.Labs {
		byID[row.ID] = row
	}
	if byID["complete-lab"].Status != progress.LabStatusComplete {
		t.Errorf("complete-lab status = %q", byID["complete-lab"].Status)
	}
	if byID["complete-lab"].PostmortemHref != "/postmortems/complete-lab" {
		t.Errorf("complete-lab href = %q", byID["complete-lab"].PostmortemHref)
	}
	if byID["active-lab"].Status != progress.LabStatusInProgress {
		t.Errorf("active-lab status = %q", byID["active-lab"].Status)
	}
	if byID["untouched-lab"].Status != progress.LabStatusAvailable {
		t.Errorf("untouched-lab status = %q", byID["untouched-lab"].Status)
	}
	if byID["untouched-lab"].Difficulty != 3 || byID["untouched-lab"].EstimatedMinutes != 30 {
		t.Errorf("untouched-lab = %+v", byID["untouched-lab"])
	}
	if byID["complete-lab"].TicketID != webcore.TicketID(0) {
		t.Errorf("complete-lab TicketID = %q, want %q", byID["complete-lab"].TicketID, webcore.TicketID(0))
	}
	if byID["active-lab"].TicketID != webcore.TicketID(1) {
		t.Errorf("active-lab TicketID = %q, want %q", byID["active-lab"].TicketID, webcore.TicketID(1))
	}
	if byID["untouched-lab"].TicketID != "" {
		t.Errorf("untouched-lab TicketID = %q, want empty (no recorded attempt)", byID["untouched-lab"].TicketID)
	}
	if byID["complete-lab"].Summary != "A one-line summary." {
		t.Errorf("complete-lab Summary = %q", byID["complete-lab"].Summary)
	}
}

func TestBuildRoadmapExactlyOneLabIsUpNext(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "complete-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
		{Lab: "active-lab", Status: progress.StatusActive},
	}}

	view := buildRoadmap(cat, pf)

	var upNext []string
	for _, stage := range view.Stages {
		for _, row := range stage.Labs {
			if row.IsUpNext {
				upNext = append(upNext, row.ID)
			}
		}
	}
	if len(upNext) != 1 || upNext[0] != "active-lab" {
		t.Errorf("IsUpNext rows = %v, want exactly [active-lab] (the first non-complete lab in roadmap order)", upNext)
	}
}

func TestBuildRoadmapUpNextSkipsLockedLab(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux", Title: "Linux fundamentals"}, Labs: []catalog.Lab{
			{ID: "lnx-001", Title: "First", Difficulty: 1},
			{ID: "lnx-002", Title: "Second", Difficulty: 1, Prerequisites: []string{"lnx-001"}},
		}},
	}}

	view := buildRoadmap(cat, &progress.ProgressFile{})

	byID := map[string]webcore.LabRowView{}
	for _, row := range view.Stages[0].Labs {
		byID[row.ID] = row
	}
	if byID["lnx-002"].Status != progress.LabStatusLocked {
		t.Fatalf("lnx-002 status = %q, want locked", byID["lnx-002"].Status)
	}
	if byID["lnx-002"].IsUpNext {
		t.Error("a locked lab must never be the Up next cursor")
	}
	if !byID["lnx-001"].IsUpNext {
		t.Error("lnx-001 (unlocked, not complete) should be the Up next cursor")
	}
}

func TestBuildRoadmapReviewDueOnlyForLabPastInterval(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "due-lab", Title: "Due lab"},
			{ID: "fresh-lab", Title: "Fresh lab"},
		}},
	}}
	old := time.Now().Add(-(progress.ReviewInterval + time.Hour))
	recent := time.Now().Add(-time.Hour)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "due-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(old)},
		{Lab: "fresh-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(recent)},
	}}

	view := buildRoadmap(cat, pf)
	byID := map[string]webcore.LabRowView{}
	for _, row := range view.Stages[0].Labs {
		byID[row.ID] = row
	}
	if !byID["due-lab"].ReviewDue {
		t.Error("due-lab solved past the review interval should be ReviewDue")
	}
	if byID["fresh-lab"].ReviewDue {
		t.Error("fresh-lab solved recently should not be ReviewDue")
	}
}

func TestBuildRoadmapCompleteLabWithoutPostmortemHasNoLink(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{{ID: "web-down"}}},
	}}
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}

	view := buildRoadmap(cat, pf)
	row := view.Stages[0].Labs[0]
	if row.Status != progress.LabStatusComplete {
		t.Errorf("status = %q, want complete", row.Status)
	}
	if row.PostmortemHref != "" {
		t.Errorf("href = %q, want empty (no postmortem recorded)", row.PostmortemHref)
	}
}

func TestBuildRoadmapShowSkip(t *testing.T) {
	cat := fixtureCatalogForRoadmap()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "complete-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
		{Lab: "active-lab", Status: progress.StatusActive},
	}}

	view := buildRoadmap(cat, pf)
	if !view.Stages[0].ShowSkip {
		t.Error("linux stage has an available and an in-progress lab, so ShowSkip should be true")
	}

	allComplete := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "git-lab", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}
	view = buildRoadmap(cat, allComplete)
	for _, stage := range view.Stages {
		if stage.ID == "git" && stage.ShowSkip {
			t.Error("git stage is fully complete, so ShowSkip should be false")
		}
	}
}
