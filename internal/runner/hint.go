package runner

import (
	"context"
	"errors"
	"fmt"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

// ErrOutOfHints is returned by Hint when every hint for the active lab
// and variant has already been shown.
var ErrOutOfHints = errors.New("no hints left for this attempt; try `breaklab reset` for a fresh attempt")

// ErrReviewHintsUnavailable is returned by Hint when the active attempt
// is a review attempt (the lab already has a complete attempt): the
// whole point is confirming unaided recall, so hints stay withheld for
// its duration.
var ErrReviewHintsUnavailable = errors.New("hints are unavailable during a review attempt — this lab is already complete; see if you can still solve it unaided")

// Hint returns the next unseen hint for the active lab: all of the lab's
// general hints, then all of the active variant's hints, in the order
// each is written in lab.yaml, incrementing the session's hint count.
// Unlike Start/Check/Reset/Stop, it needs no Docker — hints are pure
// catalog and session data.
func (r *Runner) Hint(ctx context.Context) (string, error) {
	session, err := LoadSession(r.SessionPath)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", ErrNoActiveLab
	}

	pf, err := progress.Load(r.ProgressPath)
	if err != nil {
		return "", err
	}
	if progress.IsReviewAttempt(pf, session.Lab) {
		return "", ErrReviewHintsUnavailable
	}

	lab, err := r.findLab(session.Lab)
	if err != nil {
		return "", err
	}

	var variant catalog.Variant
	found := false
	for _, v := range lab.Variants {
		if v.ID == session.Variant {
			variant, found = v, true
		}
	}
	if !found {
		return "", fmt.Errorf("variant of the active session no longer exists in lab %q; run `breaklab start %s`", lab.ID, lab.ID)
	}

	hints := append(append([]string{}, lab.Hints...), variant.Hints...)
	if session.HintsUsed >= len(hints) {
		return "", ErrOutOfHints
	}

	hint := hints[session.HintsUsed]
	session.HintsUsed++
	if err := SaveSession(r.SessionPath, *session); err != nil {
		return "", err
	}
	return hint, nil
}
