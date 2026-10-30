package progress

import "time"

// ReviewInterval is how long a lab stays "complete" before it's due for
// an unaided re-solve — long enough that genuine retention, not
// short-term memory, is what's being tested.
const ReviewInterval = 60 * 24 * time.Hour

// ReviewDue returns every lab, keyed by id, whose most recent complete
// attempt's SolvedAt is more than ReviewInterval in the past. A lab with
// no complete attempt is never due.
func ReviewDue(pf *ProgressFile) map[string]bool {
	due := make(map[string]bool)
	now := time.Now()
	for labID, attempt := range LatestCompleteAttempts(pf) {
		if now.Sub(*attempt.SolvedAt) > ReviewInterval {
			due[labID] = true
		}
	}
	return due
}

// IsReviewAttempt reports whether labID already has a complete attempt
// in pf — the condition that makes a newly started attempt on that lab
// a review attempt (an unaided re-solve) rather than a first-time solve.
func IsReviewAttempt(pf *ProgressFile, labID string) bool {
	_, found := LatestCompleteAttempt(pf, labID)
	return found
}
