package progress

import (
	"testing"
	"time"
)

func solvedPtr(t time.Time) *time.Time { return &t }

func TestLatestCompleteAttemptReturnsLatest(t *testing.T) {
	early := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	late := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Variant: "config-typo", SolvedAt: solvedPtr(early), Status: StatusComplete},
		{Lab: "web-down", Variant: "port-taken", SolvedAt: solvedPtr(late), Status: StatusComplete},
	}}

	got, ok := LatestCompleteAttempt(pf, "web-down")
	if !ok || got.Variant != "port-taken" {
		t.Errorf("got %+v, ok %v", got, ok)
	}
}

func TestLatestCompleteAttemptAbsentWhenNoneComplete(t *testing.T) {
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusSolved, SolvedAt: solvedPtr(time.Now())},
		{Lab: "web-down", Status: StatusAbandoned},
	}}

	if _, ok := LatestCompleteAttempt(pf, "web-down"); ok {
		t.Error("expected no complete attempt")
	}
}

func TestLatestCompleteAttemptsEmptyForEmptyProgress(t *testing.T) {
	got := LatestCompleteAttempts(&ProgressFile{})
	if len(got) != 0 {
		t.Errorf("got %v, want empty map", got)
	}
}

func TestLatestCompleteAttemptsKeyedByLab(t *testing.T) {
	early := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	late := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Variant: "config-typo", SolvedAt: solvedPtr(early), Status: StatusComplete},
		{Lab: "web-down", Variant: "port-taken", SolvedAt: solvedPtr(late), Status: StatusComplete},
		{Lab: "disk-full", Status: StatusActive},
	}}

	got := LatestCompleteAttempts(pf)
	if len(got) != 1 {
		t.Fatalf("got %+v, want exactly one lab", got)
	}
	if got["web-down"].Variant != "port-taken" {
		t.Errorf("web-down = %+v, want the later attempt", got["web-down"])
	}
	if _, ok := got["disk-full"]; ok {
		t.Error("disk-full has no complete attempt and should be absent")
	}
}
