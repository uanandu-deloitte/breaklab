package dashboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/runner"
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
