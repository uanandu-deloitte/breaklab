package dashboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/runner"
	"breaklab/internal/web/phases"
	"breaklab/internal/web/webcore"
)

func solvedAt(t time.Time) *time.Time { return &t }

func emptyCatalog() *catalog.Catalog { return &catalog.Catalog{} }

func TestBuildDashboardNoActiveLab(t *testing.T) {
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now()), Postmortem: "postmortems/x.md"},
	}}

	view := buildDashboard(context.Background(), emptyCatalog(), pf, nil, time.Time{}, nil, nil)
	if view.Active != nil {
		t.Errorf("Active = %+v, want nil", view.Active)
	}
	if view.LabsCompleted != 1 {
		t.Errorf("LabsCompleted = %d, want 1", view.LabsCompleted)
	}
}

func TestBuildDashboardActiveLab(t *testing.T) {
	started := time.Now().Add(-18 * time.Minute)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: started, Status: progress.StatusActive},
	}}
	session := &runner.Session{Lab: "web-down", Container: "breaklab-web-down", AttemptIndex: 0, HintsUsed: 2}

	view := buildDashboard(context.Background(), emptyCatalog(), pf, session, time.Time{}, nil, nil)
	if view.Active == nil {
		t.Fatal("Active = nil, want a block")
	}
	if view.Active.Lab != "web-down" || view.Active.HintsUsed != 2 || view.Active.ElapsedMinutes < 17 {
		t.Errorf("Active = %+v", view.Active)
	}
	if view.Active.Href != "/labs/web-down" {
		t.Errorf("Href = %q, want %q", view.Active.Href, "/labs/web-down")
	}
	if view.Active.ReminderCmd != "breaklab check" {
		t.Errorf("ReminderCmd = %q, want %q", view.Active.ReminderCmd, "breaklab check")
	}
	if view.Active.TicketID != webcore.TicketID(0) {
		t.Errorf("TicketID = %q, want %q", view.Active.TicketID, webcore.TicketID(0))
	}
}

func TestBuildDashboardDockerTelemetryDegradesOnError(t *testing.T) {
	pf := &progress.ProgressFile{}
	failingVersion := webcore.VersionFunc(func(ctx context.Context) (string, error) {
		return "", errors.New("docker unreachable")
	})
	view := buildDashboard(context.Background(), emptyCatalog(), pf, nil, time.Time{}, failingVersion, nil)
	if view.DockerVersion != "" {
		t.Errorf("DockerVersion = %q, want empty on error", view.DockerVersion)
	}
	if view.Uptime != "" {
		t.Errorf("Uptime = %q, want empty for a zero start time", view.Uptime)
	}
}

func TestBuildDashboardPatternsRow(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "k8s"}, Labs: []catalog.Lab{
			{ID: "net-1", Tags: []string{"networking"}},
			{ID: "net-2", Tags: []string{"networking"}},
			{ID: "net-3", Tags: []string{"networking"}},
			{ID: "lnx-1", Tags: []string{"linux"}},
			{ID: "lnx-2", Tags: []string{"linux"}},
			{ID: "lnx-3", Tags: []string{"linux"}},
		}},
	}}
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "net-1", Status: progress.StatusComplete, HintsUsed: 6},
		{Lab: "net-2", Status: progress.StatusComplete, HintsUsed: 6},
		{Lab: "net-3", Status: progress.StatusComplete, HintsUsed: 6},
		{Lab: "lnx-1", Status: progress.StatusComplete, HintsUsed: 1},
		{Lab: "lnx-2", Status: progress.StatusComplete, HintsUsed: 1},
		{Lab: "lnx-3", Status: progress.StatusComplete, HintsUsed: 1},
	}}

	view := buildDashboard(context.Background(), cat, pf, nil, time.Time{}, nil, nil)
	if len(view.Patterns) != 1 || view.Patterns[0].Tag != "networking" {
		t.Fatalf("Patterns = %+v, want exactly one pattern for networking", view.Patterns)
	}
}

func TestBuildDashboardNoPatternsWhenNoneQualify(t *testing.T) {
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Status: progress.StatusComplete, HintsUsed: 1},
	}}
	view := buildDashboard(context.Background(), emptyCatalog(), pf, nil, time.Time{}, nil, nil)
	if len(view.Patterns) != 0 {
		t.Errorf("Patterns = %+v, want empty", view.Patterns)
	}
}

// add-curriculum-phases task 3.1: stage cards use the same
// complete-attempt definition LabsCompleted already uses, and a stage
// isn't Complete until every one of its labs has a complete attempt.
func TestBuildDashboardStageCards(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux", Phase: "phase-1", Title: "Linux fundamentals"}, Labs: []catalog.Lab{
			{ID: "web-down"}, {ID: "disk-full"},
		}},
		{Stage: catalog.Stage{ID: "git", Phase: "phase-1", Title: "Git"}, Labs: []catalog.Lab{
			{ID: "bad-merge"},
		}},
		{Stage: catalog.Stage{ID: "capstone", Phase: "phase-1", Title: "Capstone", Capstone: true}, Labs: []catalog.Lab{
			{ID: "final-incident"},
		}},
	}}
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
		{Lab: "bad-merge", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}

	view := buildDashboard(context.Background(), cat, pf, nil, time.Time{}, nil, nil)
	if len(view.Stages) != 2 {
		t.Fatalf("Stages = %+v, want 2 (capstone excluded)", view.Stages)
	}

	byID := map[string]StageCardView{}
	for _, c := range view.Stages {
		byID[c.ID] = c
	}
	linux := byID["linux"]
	if linux.PhaseID != "phase-1" || linux.LabsCompleted != 1 || linux.LabsTotal != 2 || linux.Complete {
		t.Errorf("linux card = %+v, want 1/2 complete, not Complete", linux)
	}
	git := byID["git"]
	if git.PhaseID != "phase-1" || git.LabsCompleted != 1 || git.LabsTotal != 1 || !git.Complete {
		t.Errorf("git card = %+v, want 1/1 complete, Complete", git)
	}
}

// add-curriculum-phases task 3.2: DashboardView.Phases matches a direct
// call to phases.BuildPhaseList for the same catalog/progress, so the
// dashboard never computes "phase progress" a second, different way.
func TestBuildDashboardPhasesMatchesPhaseListBuilder(t *testing.T) {
	cat := &catalog.Catalog{
		Stages: []catalog.StageGroup{
			{Stage: catalog.Stage{ID: "linux", Phase: "phase-1"}, Labs: []catalog.Lab{{ID: "web-down"}}},
		},
		Phases: []catalog.PhaseGroup{
			{Phase: catalog.Phase{ID: "phase-1", Title: "Phase 1"}, Stages: []catalog.StageGroup{
				{Stage: catalog.Stage{ID: "linux", Phase: "phase-1"}, Labs: []catalog.Lab{{ID: "web-down"}}},
			}},
			{Phase: catalog.Phase{ID: "phase-2", Title: "Phase 2 — Cloud"}, Stages: nil},
		},
	}
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Status: progress.StatusComplete, SolvedAt: solvedAt(time.Now())},
	}}

	view := buildDashboard(context.Background(), cat, pf, nil, time.Time{}, nil, nil)
	want := phases.BuildPhaseList(cat, pf).Phases
	if len(view.Phases) != len(want) {
		t.Fatalf("Phases = %+v, want %+v", view.Phases, want)
	}
	for i := range want {
		if view.Phases[i] != want[i] {
			t.Errorf("Phases[%d] = %+v, want %+v", i, view.Phases[i], want[i])
		}
	}
	if len(view.Phases) != 2 || view.Phases[1].ID != "phase-2" {
		t.Fatalf("expected phase-2 to still appear with zero stages, got %+v", view.Phases)
	}
}

// add-curriculum-phases task 3.3: attempts in the current month are
// counted per day, with a deterministic most-engaged day on a tie.
func TestBuildDashboardMonthlyEngagement(t *testing.T) {
	now := time.Date(2026, 11, 12, 15, 0, 0, 0, time.UTC)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "a", StartedAt: time.Date(2026, 11, 10, 9, 0, 0, 0, time.UTC)},
		{Lab: "b", StartedAt: time.Date(2026, 11, 10, 10, 0, 0, 0, time.UTC)},
		{Lab: "c", StartedAt: time.Date(2026, 11, 11, 9, 0, 0, 0, time.UTC)},
		// Outside the month: must not be counted.
		{Lab: "d", StartedAt: time.Date(2026, 10, 30, 9, 0, 0, 0, time.UTC)},
	}}

	counts, mostEngagedDay := buildMonthlyEngagement(pf.Attempts, now)
	if counts["2026-11-10"] != 2 || counts["2026-11-11"] != 1 {
		t.Fatalf("counts = %+v, want 2026-11-10:2, 2026-11-11:1", counts)
	}
	if len(counts) != 2 {
		t.Fatalf("counts = %+v, want exactly 2 days (out-of-month attempt excluded)", counts)
	}
	if mostEngagedDay == nil || *mostEngagedDay != "2026-11-10" {
		t.Fatalf("mostEngagedDay = %v, want 2026-11-10", mostEngagedDay)
	}
}

// Monthly engagement with no activity this month still renders an
// empty (non-nil) map, per the web-ui spec's scenario of the same name.
func TestBuildDashboardMonthlyEngagementNoActivity(t *testing.T) {
	view := buildDashboard(context.Background(), emptyCatalog(), &progress.ProgressFile{}, nil, time.Time{}, nil, nil)
	if view.MonthlyEngagement == nil {
		t.Fatal("MonthlyEngagement = nil, want a non-nil empty map")
	}
	if len(view.MonthlyEngagement) != 0 {
		t.Fatalf("MonthlyEngagement = %+v, want empty", view.MonthlyEngagement)
	}
	if view.MostEngagedDay != nil {
		t.Fatalf("MostEngagedDay = %v, want nil", view.MostEngagedDay)
	}
}

func TestBuildDashboardRecentCappedAtFive(t *testing.T) {
	pf := &progress.ProgressFile{}
	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		pf.Attempts = append(pf.Attempts, progress.Attempt{
			Lab:        "lab-" + string(rune('a'+i)),
			Variant:    "v",
			Status:     progress.StatusComplete,
			StartedAt:  base.Add(time.Duration(i) * 24 * time.Hour),
			SolvedAt:   solvedAt(base.Add(time.Duration(i)*24*time.Hour + time.Hour)),
			Postmortem: "postmortems/lab-" + string(rune('a'+i)) + ".md",
		})
	}

	view := buildDashboard(context.Background(), emptyCatalog(), pf, nil, time.Time{}, nil, nil)
	if len(view.Recent) != 5 {
		t.Fatalf("Recent has %d entries, want 5", len(view.Recent))
	}
	// Most recent (lab-g, added last with the latest SolvedAt) first.
	if view.Recent[0].Lab != "lab-g" {
		t.Errorf("Recent[0].Lab = %q, want lab-g (most recently solved)", view.Recent[0].Lab)
	}
	if view.Recent[0].TicketID != webcore.TicketID(6) {
		t.Errorf("Recent[0].TicketID = %q, want %q (attempt index 6)", view.Recent[0].TicketID, webcore.TicketID(6))
	}
}
