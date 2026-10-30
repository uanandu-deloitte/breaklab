package progress

import (
	"testing"
	"time"
)

func solvedAttempt(lab string, solvedAt time.Time) Attempt {
	t := solvedAt
	return Attempt{Lab: lab, Status: StatusSolved, SolvedAt: &t}
}

func TestStreakNoAttempts(t *testing.T) {
	now := time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC)
	if got := Streak(&ProgressFile{}, now); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestStreakSolvedOnlyToday(t *testing.T) {
	now := time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		solvedAttempt("web-down", time.Date(2026, 10, 17, 8, 0, 0, 0, time.UTC)),
	}}
	if got := Streak(pf, now); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestStreakSolvedOnlyYesterday(t *testing.T) {
	now := time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		solvedAttempt("web-down", time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)),
	}}
	if got := Streak(pf, now); got != 1 {
		t.Errorf("got %d, want 1 (grace day for today)", got)
	}
}

func TestStreakThreeConsecutiveDaysEndingYesterday(t *testing.T) {
	now := time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		solvedAttempt("a", time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)),
		solvedAttempt("b", time.Date(2026, 10, 15, 8, 0, 0, 0, time.UTC)),
		solvedAttempt("c", time.Date(2026, 10, 14, 8, 0, 0, 0, time.UTC)),
	}}
	if got := Streak(pf, now); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestStreakGapStopsCount(t *testing.T) {
	now := time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		solvedAttempt("a", time.Date(2026, 10, 17, 8, 0, 0, 0, time.UTC)), // today
		solvedAttempt("b", time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)), // yesterday
		// 2026-10-15 missing: gap
		solvedAttempt("c", time.Date(2026, 10, 14, 8, 0, 0, 0, time.UTC)),
	}}
	if got := Streak(pf, now); got != 2 {
		t.Errorf("got %d, want 2 (gap two days back stops the count)", got)
	}
}

func TestStreakNoActivityAtAllGraceDayAlsoEmpty(t *testing.T) {
	now := time.Date(2026, 10, 17, 9, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		solvedAttempt("a", time.Date(2026, 10, 10, 8, 0, 0, 0, time.UTC)),
	}}
	if got := Streak(pf, now); got != 0 {
		t.Errorf("got %d, want 0 (last activity too long ago, grace day also empty)", got)
	}
}
