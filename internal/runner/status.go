package runner

import (
	"context"
	"fmt"
	"sort"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

// ActiveStatus is the active lab's block in `breaklab status`, or nil
// when no lab is running.
type ActiveStatus struct {
	Lab            string
	ElapsedMinutes int
	HintsUsed      int
}

// StatusView is everything `breaklab status` prints.
type StatusView struct {
	Stages    []progress.StageStatus
	Capstone  *progress.CapstoneStatus // nil if roadmap.yaml declares no capstone stage
	Streak    int
	Active    *ActiveStatus
	ReviewDue []string           // lab ids due for review, sorted; empty if none
	Patterns  []progress.Pattern // topic tags needing consistently more hints; empty if none
}

// Status assembles the view for `breaklab status`. Like Hint, it needs
// no Docker — it's pure catalog, progress, and session data.
func (r *Runner) Status(ctx context.Context) (StatusView, error) {
	cat, err := catalog.LoadCatalog(r.LabsDir, r.RoadmapPath)
	if err != nil {
		return StatusView{}, err
	}
	pf, err := progress.Load(r.ProgressPath)
	if err != nil {
		return StatusView{}, err
	}

	view := StatusView{
		Stages:    progress.StageCompletion(cat, pf),
		Capstone:  progress.Capstone(cat, pf),
		Streak:    progress.Streak(pf, time.Now()),
		ReviewDue: sortedKeys(progress.ReviewDue(pf)),
		Patterns:  progress.Patterns(cat, pf),
	}

	session, err := LoadSession(r.SessionPath)
	if err != nil {
		return StatusView{}, err
	}
	if session == nil {
		return view, nil
	}
	if session.AttemptIndex < 0 || session.AttemptIndex >= len(pf.Attempts) {
		return StatusView{}, fmt.Errorf("session references unknown attempt %d in %s", session.AttemptIndex, r.ProgressPath)
	}

	// Elapsed time comes from the attempt's own started_at, not the
	// session, so a `breaklab reset` mid-attempt doesn't reset the clock
	// (design.md Decision 2). Hints come live from the session, since
	// the persisted record isn't updated until the attempt is finalized
	// (design.md Decision 5).
	attempt := pf.Attempts[session.AttemptIndex]
	view.Active = &ActiveStatus{
		Lab:            session.Lab,
		ElapsedMinutes: int(time.Since(attempt.StartedAt).Minutes()),
		HintsUsed:      session.HintsUsed,
	}
	return view, nil
}

// sortedKeys returns m's keys in sorted order, for a deterministic
// listing rather than Go's randomized map iteration order.
func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
