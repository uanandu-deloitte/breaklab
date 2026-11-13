package progress

// LatestCompleteAttempt returns labID's most recently solved attempt
// with status "complete", if any.
func LatestCompleteAttempt(pf *ProgressFile, labID string) (Attempt, bool) {
	var best Attempt
	found := false
	for _, a := range pf.Attempts {
		if a.Lab != labID || a.Status != StatusComplete || a.SolvedAt == nil {
			continue
		}
		if !found || a.SolvedAt.After(*best.SolvedAt) {
			best, found = a, true
		}
	}
	return best, found
}

// LatestCompleteAttempts returns, for every lab with at least one
// "complete" attempt, that lab's most recently solved one, keyed by lab id.
func LatestCompleteAttempts(pf *ProgressFile) map[string]Attempt {
	latest := make(map[string]Attempt)
	for _, a := range pf.Attempts {
		if a.Status != StatusComplete || a.SolvedAt == nil {
			continue
		}
		if existing, ok := latest[a.Lab]; !ok || a.SolvedAt.After(*existing.SolvedAt) {
			latest[a.Lab] = a
		}
	}
	return latest
}

// LatestAttemptIndex returns the index within pf.Attempts of labID's
// most recently started attempt of any status, if any.
func LatestAttemptIndex(pf *ProgressFile, labID string) (int, bool) {
	best := -1
	for i, a := range pf.Attempts {
		if a.Lab != labID {
			continue
		}
		if best == -1 || a.StartedAt.After(pf.Attempts[best].StartedAt) {
			best = i
		}
	}
	return best, best != -1
}

// LatestAttempt returns labID's most recently started attempt of any
// status, if any.
func LatestAttempt(pf *ProgressFile, labID string) (Attempt, bool) {
	i, ok := LatestAttemptIndex(pf, labID)
	if !ok {
		return Attempt{}, false
	}
	return pf.Attempts[i], true
}
