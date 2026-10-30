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

func TestMergeSyncedAddsOnlyUnknownHostedAttempts(t *testing.T) {
	pf := &ProgressFile{Version: 1, Attempts: []Attempt{
		{ID: "local-1", Lab: "web-down", Status: StatusSolved},
	}}
	hosted := []Attempt{
		{ID: "local-1", Lab: "web-down", Status: StatusComplete}, // already known locally — must not duplicate
		{ID: "hosted-1", Lab: "stuck-lock", Status: StatusSolved},
	}

	merged := MergeSynced(pf, hosted)
	if len(merged.Attempts) != 2 {
		t.Fatalf("merged attempts = %+v, want exactly 2", merged.Attempts)
	}
	if merged.Attempts[0].ID != "local-1" || merged.Attempts[0].Status != StatusSolved {
		t.Errorf("local attempt should be untouched (local wins), got %+v", merged.Attempts[0])
	}
	if merged.Attempts[1].ID != "hosted-1" {
		t.Errorf("hosted-only attempt should be added, got %+v", merged.Attempts[1])
	}
}

func TestMergeSyncedNeverMutatesInput(t *testing.T) {
	pf := &ProgressFile{Version: 1, Attempts: []Attempt{{ID: "local-1", Lab: "web-down"}}}

	MergeSynced(pf, []Attempt{{ID: "hosted-1", Lab: "stuck-lock"}})

	if len(pf.Attempts) != 1 {
		t.Errorf("pf.Attempts = %+v, want the original untouched (MergeSynced must not mutate its input)", pf.Attempts)
	}
}

func TestMergeSyncedSkipsHostedAttemptsWithNoID(t *testing.T) {
	pf := &ProgressFile{Version: 1}
	hosted := []Attempt{{ID: "", Lab: "web-down"}}

	merged := MergeSynced(pf, hosted)
	if len(merged.Attempts) != 0 {
		t.Errorf("merged = %+v, want empty: a hosted attempt with no ID can never be matched/trusted", merged.Attempts)
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
