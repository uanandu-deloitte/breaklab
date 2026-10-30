package progress

import (
	"testing"
	"time"
)

func TestReviewDueNeverCompletedLab(t *testing.T) {
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusActive},
	}}
	if due := ReviewDue(pf); due["web-down"] {
		t.Error("a lab with no complete attempt must never be review-due")
	}
}

func TestReviewDueRecentlyCompletedLabNotDue(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour)
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusComplete, SolvedAt: &recent},
	}}
	if due := ReviewDue(pf); due["web-down"] {
		t.Error("a lab solved well within the interval must not be review-due")
	}
}

func TestReviewDueLabPastIntervalIsDue(t *testing.T) {
	old := time.Now().Add(-(ReviewInterval + time.Hour))
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusComplete, SolvedAt: &old},
	}}
	if due := ReviewDue(pf); !due["web-down"] {
		t.Error("a lab solved past the review interval must be review-due")
	}
}

func TestIsReviewAttemptTrueWhenLabAlreadyComplete(t *testing.T) {
	solved := time.Now().Add(-24 * time.Hour)
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusComplete, SolvedAt: &solved},
	}}
	if !IsReviewAttempt(pf, "web-down") {
		t.Error("expected a review attempt: web-down already has a complete attempt")
	}
}

func TestIsReviewAttemptFalseForFirstTimeSolve(t *testing.T) {
	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusActive},
	}}
	if IsReviewAttempt(pf, "web-down") {
		t.Error("expected not a review attempt: web-down has no prior complete attempt")
	}
}
