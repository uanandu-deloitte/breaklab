package progress

import "breaklab/internal/catalog"

// StageStatus is one roadmap stage's completion count for `breaklab status`.
type StageStatus struct {
	ID        string
	Title     string
	Completed int
	Total     int
}

// completedLabsSet returns the set of lab ids with at least one attempt
// in pf whose status is "solved", "complete", or "skipped".
func completedLabsSet(pf *ProgressFile) map[string]bool {
	completed := make(map[string]bool)
	for _, a := range pf.Attempts {
		if a.Status == StatusSolved || a.Status == StatusComplete || a.Status == StatusSkipped {
			completed[a.Lab] = true
		}
	}
	return completed
}

// StageCompletion returns one StageStatus per regular (non-capstone)
// roadmap stage, in roadmap order. A lab counts as completed if any of
// its attempts in pf has status "solved" or "complete".
func StageCompletion(cat *catalog.Catalog, pf *ProgressFile) []StageStatus {
	completedLabs := completedLabsSet(pf)

	regular := cat.RegularStages()
	result := make([]StageStatus, 0, len(regular))
	for _, group := range regular {
		status := StageStatus{ID: group.Stage.ID, Title: group.Stage.Title, Total: len(group.Labs)}
		for _, lab := range group.Labs {
			if completedLabs[lab.ID] {
				status.Completed++
			}
		}
		result = append(result, status)
	}
	return result
}

// CapstoneStatus is the capstone stage's own completion/lock state.
type CapstoneStatus struct {
	ID        string
	Title     string
	Completed int
	Total     int

	// Unlocked is true only once every regular (non-capstone) stage is
	// fully complete — the capstone is an earned finale, not just
	// another entry in the list.
	Unlocked bool
}

// Capstone returns the capstone stage's status, or nil if roadmap.yaml
// declares no capstone stage.
func Capstone(cat *catalog.Catalog, pf *ProgressFile) *CapstoneStatus {
	group := cat.CapstoneStage()
	if group == nil {
		return nil
	}

	unlocked := true
	for _, s := range StageCompletion(cat, pf) {
		if s.Completed < s.Total {
			unlocked = false
			break
		}
	}

	completedLabs := completedLabsSet(pf)
	status := &CapstoneStatus{ID: group.Stage.ID, Title: group.Stage.Title, Total: len(group.Labs), Unlocked: unlocked}
	for _, lab := range group.Labs {
		if completedLabs[lab.ID] {
			status.Completed++
		}
	}
	return status
}
