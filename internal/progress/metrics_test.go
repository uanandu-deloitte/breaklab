package progress

import (
	"testing"
	"time"

	"breaklab/internal/catalog"
)

func solvedAt(t time.Time) *time.Time { return &t }

func TestAvgResolutionTime(t *testing.T) {
	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	attempts := []Attempt{
		{Status: StatusComplete, StartedAt: base, SolvedAt: solvedAt(base.Add(10 * time.Minute))},
		{Status: StatusSolved, StartedAt: base, SolvedAt: solvedAt(base.Add(20 * time.Minute))},
		{Status: StatusActive, StartedAt: base}, // no SolvedAt: excluded
		{Status: StatusAbandoned, StartedAt: base, SolvedAt: solvedAt(base.Add(time.Hour))},
	}
	got := AvgResolutionTime(attempts)
	want := 15 * time.Minute
	if got != want {
		t.Errorf("AvgResolutionTime = %v, want %v", got, want)
	}
}

func TestAvgResolutionTimeNoAttempts(t *testing.T) {
	if got := AvgResolutionTime(nil); got != 0 {
		t.Errorf("AvgResolutionTime(nil) = %v, want 0", got)
	}
}

func TestPostmortemCounts(t *testing.T) {
	attempts := []Attempt{
		{Status: StatusComplete},
		{Status: StatusComplete},
		{Status: StatusSolved},
		{Status: StatusActive},
		{Status: StatusAbandoned},
	}
	drafts, published := PostmortemCounts(attempts)
	if drafts != 1 || published != 2 {
		t.Errorf("PostmortemCounts = (%d, %d), want (1, 2)", drafts, published)
	}
}

func testCatalogForMetrics() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{
			Stage: catalog.Stage{ID: "linux", Title: "Linux"},
			Labs: []catalog.Lab{
				{ID: "web-down", Difficulty: 1},
				{ID: "zombie-procs", Difficulty: 1},
				{ID: "dns-cascade", Difficulty: 3},
			},
		},
	}}
}

func TestBuildTrainingVelocity(t *testing.T) {
	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusComplete, StartedAt: base, SolvedAt: solvedAt(base.Add(10 * time.Minute)), Postmortem: "x.md"},
		{Lab: "zombie-procs", Status: StatusComplete, StartedAt: base, SolvedAt: solvedAt(base.Add(20 * time.Minute)), Postmortem: "y.md"},
	}}

	tv := BuildTrainingVelocity(testCatalogForMetrics(), pf)

	if tv.SolvedCount != 2 || tv.TotalCount != 3 {
		t.Errorf("SolvedCount/TotalCount = %d/%d, want 2/3", tv.SolvedCount, tv.TotalCount)
	}
	wantPercent := 100.0 * 2 / 3
	if tv.PercentSolved != wantPercent {
		t.Errorf("PercentSolved = %v, want %v", tv.PercentSolved, wantPercent)
	}
	if len(tv.DifficultyCurve) != 2 {
		t.Fatalf("DifficultyCurve has %d buckets, want 2", len(tv.DifficultyCurve))
	}
	if tv.DifficultyCurve[0].Difficulty != 1 || tv.DifficultyCurve[0].Count != 2 || tv.DifficultyCurve[0].Solved != 2 {
		t.Errorf("DifficultyCurve[0] = %+v, want {1 2 2}", tv.DifficultyCurve[0])
	}
	if tv.DifficultyCurve[1].Difficulty != 3 || tv.DifficultyCurve[1].Count != 1 || tv.DifficultyCurve[1].Solved != 0 {
		t.Errorf("DifficultyCurve[1] = %+v, want {3 1 0}", tv.DifficultyCurve[1])
	}
	wantMedian := 15 * time.Minute
	if tv.MedianSolveTime != wantMedian {
		t.Errorf("MedianSolveTime = %v, want %v", tv.MedianSolveTime, wantMedian)
	}
}

func TestBuildTrainingVelocityNoAttempts(t *testing.T) {
	tv := BuildTrainingVelocity(testCatalogForMetrics(), &ProgressFile{})
	if tv.SolvedCount != 0 || tv.PercentSolved != 0 {
		t.Errorf("expected zero solved/percent, got %+v", tv)
	}
	if tv.MedianSolveTime != 0 {
		t.Errorf("MedianSolveTime = %v, want 0", tv.MedianSolveTime)
	}
}
