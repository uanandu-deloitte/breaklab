package site

import (
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

func solvedAt(t time.Time) *time.Time { return &t }

func fixtureCatalog() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{
			Stage: catalog.Stage{ID: "linux", Title: "Linux fundamentals"},
			Labs: []catalog.Lab{
				{ID: "web-down", Title: "The website returns nothing", Difficulty: 1},
				{ID: "disk-full", Title: "Disk is full", Difficulty: 2},
			},
		},
		{
			Stage: catalog.Stage{ID: "git", Title: "Git"},
			Labs: []catalog.Lab{
				{ID: "git-lost-work", Title: "Lost commit", Difficulty: 2},
			},
		},
	}}
}

func TestBuildIndexCountsAndRows(t *testing.T) {
	cat := fixtureCatalog()
	started := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	solved := started.Add(18 * time.Minute)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: started, SolvedAt: solvedAt(solved), Status: progress.StatusComplete, Postmortem: "postmortems/2026-10-17-web-down-config-typo.md"},
		{Lab: "disk-full", Variant: "runaway-log", StartedAt: started, Status: progress.StatusActive},
		{Lab: "git-lost-work", Variant: "hard-reset", StartedAt: started, Status: progress.StatusAbandoned},
	}}

	view := BuildIndex(cat, pf, map[string]bool{})

	if len(view.Stages) != 2 {
		t.Fatalf("stages = %+v", view.Stages)
	}
	linux := view.Stages[0]
	if linux.Completed != 1 || linux.Total != 2 || linux.Percent != 50 {
		t.Errorf("linux stage = %+v", linux)
	}
	if len(linux.Labs) != 1 {
		t.Fatalf("linux.Labs = %+v", linux.Labs)
	}
	row := linux.Labs[0]
	if row.ID != "web-down" || row.Date != "2026-10-17" || row.Difficulty != 1 || row.MinutesToFix != 18 {
		t.Errorf("row = %+v", row)
	}
	if row.PostmortemHref != "" {
		t.Errorf("PostmortemHref = %q, want empty (not in published map)", row.PostmortemHref)
	}

	git := view.Stages[1]
	if git.Completed != 0 || git.Total != 1 {
		t.Errorf("git stage = %+v", git)
	}

	if view.Totals.LabsCompleted != 1 {
		t.Errorf("LabsCompleted = %d, want 1", view.Totals.LabsCompleted)
	}
	if view.Totals.VariantsSeen != 3 {
		t.Errorf("VariantsSeen = %d, want 3 (one per attempt, any status)", view.Totals.VariantsSeen)
	}
}

func TestBuildIndexSetsHrefOnlyWhenPublished(t *testing.T) {
	cat := fixtureCatalog()
	started := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: started, SolvedAt: solvedAt(started.Add(time.Minute)), Status: progress.StatusComplete, Postmortem: "postmortems/x.md"},
	}}

	withoutPublish := BuildIndex(cat, pf, map[string]bool{})
	if withoutPublish.Stages[0].Labs[0].PostmortemHref != "" {
		t.Error("expected no href when postmortem is not published")
	}

	withPublish := BuildIndex(cat, pf, map[string]bool{"postmortems/x.md": true})
	got := withPublish.Stages[0].Labs[0]
	if got.PostmortemHref != "postmortems/x.html" {
		t.Errorf("href = %q, want postmortems/x.html", got.PostmortemHref)
	}
	if got.Date == "" || got.Difficulty == 0 {
		t.Errorf("row still needs date/difficulty regardless of publish state: %+v", got)
	}
}

func TestBuildIndexDedupesMultipleCompleteAttemptsPerLab(t *testing.T) {
	cat := fixtureCatalog()
	early := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	late := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: early, SolvedAt: solvedAt(early.Add(5 * time.Minute)), Status: progress.StatusComplete},
		{Lab: "web-down", Variant: "port-taken", StartedAt: late, SolvedAt: solvedAt(late.Add(10 * time.Minute)), Status: progress.StatusComplete},
	}}

	view := BuildIndex(cat, pf, map[string]bool{})
	if view.Totals.LabsCompleted != 1 {
		t.Errorf("LabsCompleted = %d, want 1 (deduped per lab)", view.Totals.LabsCompleted)
	}
	if len(view.Stages[0].Labs) != 1 {
		t.Fatalf("linux.Labs = %+v, want exactly one row", view.Stages[0].Labs)
	}
	if got := view.Stages[0].Labs[0].Date; got != "2026-10-17" {
		t.Errorf("date = %q, want the latest attempt's date", got)
	}
}

func TestBuildPostmortemsIndexOnlyIncludesPublished(t *testing.T) {
	cat := fixtureCatalog()
	started := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: started, SolvedAt: solvedAt(started.Add(18 * time.Minute)), Status: progress.StatusComplete, Postmortem: "postmortems/web-down.md"},
		{Lab: "disk-full", Variant: "runaway-log", StartedAt: started, SolvedAt: solvedAt(started.Add(-time.Hour)), Status: progress.StatusComplete, Postmortem: "postmortems/disk-full.md"},
	}}

	published := map[string]bool{"postmortems/web-down.md": true}
	view := BuildPostmortemsIndex(cat, pf, published)
	if len(view.Postmortems) != 1 {
		t.Fatalf("Postmortems = %+v, want exactly the one published entry", view.Postmortems)
	}
	row := view.Postmortems[0]
	if row.Title != "The website returns nothing" || row.StageTitle != "Linux fundamentals" || row.Difficulty != 1 || row.Href != "web-down.html" {
		t.Errorf("row = %+v", row)
	}
}

func TestBuildPostmortemsIndexSortedMostRecentFirst(t *testing.T) {
	cat := fixtureCatalog()
	early := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	late := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: early, SolvedAt: solvedAt(early.Add(5 * time.Minute)), Status: progress.StatusComplete, Postmortem: "postmortems/web-down.md"},
		{Lab: "disk-full", Variant: "runaway-log", StartedAt: late, SolvedAt: solvedAt(late.Add(10 * time.Minute)), Status: progress.StatusComplete, Postmortem: "postmortems/disk-full.md"},
	}}
	published := map[string]bool{"postmortems/web-down.md": true, "postmortems/disk-full.md": true}

	view := BuildPostmortemsIndex(cat, pf, published)
	if len(view.Postmortems) != 2 {
		t.Fatalf("Postmortems = %+v, want 2", view.Postmortems)
	}
	if view.Postmortems[0].Title != "Disk is full" || view.Postmortems[1].Title != "The website returns nothing" {
		t.Errorf("order = %+v, want disk-full (later) before web-down (earlier)", view.Postmortems)
	}
}

func TestBuildPostmortemsIndexEmptyWhenNonePublished(t *testing.T) {
	cat := fixtureCatalog()
	started := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: started, SolvedAt: solvedAt(started.Add(18 * time.Minute)), Status: progress.StatusComplete, Postmortem: "postmortems/web-down.md"},
	}}

	view := BuildPostmortemsIndex(cat, pf, map[string]bool{})
	if len(view.Postmortems) != 0 {
		t.Errorf("Postmortems = %+v, want empty when nothing is published", view.Postmortems)
	}
}
