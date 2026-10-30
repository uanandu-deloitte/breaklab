package progress

// MergeSynced combines pf's attempts with hosted (a hosted account's
// synced attempts, from account.Store.PullAttempts) by ID: a hosted
// attempt whose ID already exists locally is dropped — local always
// wins, since it's at least as fresh as whatever's already been pushed
// — and a hosted attempt with no local counterpart is added. An
// attempt with no ID (recorded before sync existed) is never matched
// against hosted and always passes through unchanged. Returns a new
// ProgressFile; pf itself is never mutated, and this never writes
// progress.json — only `breaklab sync` does that.
func MergeSynced(pf *ProgressFile, hosted []Attempt) *ProgressFile {
	local := make(map[string]bool, len(pf.Attempts))
	for _, a := range pf.Attempts {
		if a.ID != "" {
			local[a.ID] = true
		}
	}

	merged := append([]Attempt(nil), pf.Attempts...)
	for _, a := range hosted {
		if a.ID == "" || local[a.ID] {
			continue
		}
		merged = append(merged, a)
	}
	return &ProgressFile{Version: pf.Version, Attempts: merged}
}

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
