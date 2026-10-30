package runner

import (
	"context"

	"breaklab/internal/progress"
)

// SetHypothesis records the learner's own guess at the cause on the
// active attempt, for later comparison against the postmortem's Root
// cause section. Like Hint and Status, it needs no Docker. It's a no-op
// (per progress.SetHypothesis) if the active session's attempt has
// already been solved or abandoned out from under it.
func (r *Runner) SetHypothesis(ctx context.Context, text string) error {
	session, err := LoadSession(r.SessionPath)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrNoActiveLab
	}
	return progress.SetHypothesis(r.ProgressPath, session.AttemptIndex, text)
}
