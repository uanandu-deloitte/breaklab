package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"breaklab/internal/catalog"
	"breaklab/internal/postmortem"
	"breaklab/internal/progress"
)

// findEligibleAttempt returns the index of the latest (by SolvedAt)
// attempt with status "solved" or "complete" and no postmortem yet
// (design.md Decision 1), or ok=false if there is none.
func findEligibleAttempt(pf *progress.ProgressFile) (index int, ok bool) {
	best := -1
	for i, a := range pf.Attempts {
		if a.Status != progress.StatusSolved && a.Status != progress.StatusComplete {
			continue
		}
		if a.Postmortem != "" || a.SolvedAt == nil {
			continue
		}
		if best == -1 || a.SolvedAt.After(*pf.Attempts[best].SolvedAt) {
			best = i
		}
	}
	return best, best != -1
}

// latestPostmortemAttempt returns the index of the latest (by SolvedAt)
// attempt with a non-empty Postmortem, regardless of status (design.md
// Decision 4), or ok=false if none has one.
func latestPostmortemAttempt(pf *progress.ProgressFile) (index int, ok bool) {
	best := -1
	for i, a := range pf.Attempts {
		if a.Postmortem == "" || a.SolvedAt == nil {
			continue
		}
		if best == -1 || a.SolvedAt.After(*pf.Attempts[best].SolvedAt) {
			best = i
		}
	}
	return best, best != -1
}

// Postmortem creates a postmortem for the most recently solved attempt
// that doesn't have one yet, writes it to postmortems/, records its path
// on the attempt, and opens it in the configured editor. It needs no
// Docker, like Hint and Status.
func (r *Runner) Postmortem(ctx context.Context) (string, error) {
	pf, err := progress.Load(r.ProgressPath)
	if err != nil {
		return "", err
	}
	index, ok := findEligibleAttempt(pf)
	if !ok {
		return "", errors.New("Solve a lab first")
	}
	attempt := pf.Attempts[index]

	lab, err := r.findLab(attempt.Lab)
	if err != nil {
		return "", err
	}

	// progress.json is hand-editable; never build a filename from ids that
	// could carry a path separator or "..". Variant may legitimately be
	// empty (a lab with no variants), so it's only checked when present.
	if !catalog.ValidID(attempt.Lab) || (attempt.Variant != "" && !catalog.ValidID(attempt.Variant)) {
		return "", fmt.Errorf("progress.json attempt has an invalid lab %q or variant %q", attempt.Lab, attempt.Variant)
	}

	date := attempt.SolvedAt.Local().Format("2006-01-02")
	path, err := r.newPostmortemPath(date, attempt.Lab, attempt.Variant)
	if err != nil {
		return "", err
	}

	minutesToFix := int(attempt.SolvedAt.Sub(attempt.StartedAt).Minutes())
	fm := postmortem.FrontMatter{
		Lab:          attempt.Lab,
		Variant:      attempt.Variant,
		Stage:        lab.Stage,
		Date:         date,
		MinutesToFix: minutesToFix,
		HintsUsed:    attempt.HintsUsed,
		Publish:      false,
	}
	content, err := postmortem.Render(r.TemplatePath, fm, lab.Title, attempt.Hypothesis)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(r.PostmortemsDir, 0o755); err != nil {
		return "", fmt.Errorf("creating %s: %w", r.PostmortemsDir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", err
	}

	if err := progress.SetPostmortem(r.ProgressPath, index, path); err != nil {
		return "", err
	}
	if err := r.OpenEditor(ctx, path); err != nil {
		return "", err
	}
	return path, nil
}

// ValidationResult is the outcome of ValidatePostmortem.
type ValidationResult struct {
	Valid   bool
	Failing []string

	// RootCauseWarning is non-empty when the Root cause section shares
	// no meaningful word with the solved attempt's actual variant's
	// hints — an advisory nudge only. It never affects Valid/Failing or
	// whether the attempt gets promoted to "complete".
	RootCauseWarning string
}

// ValidatePostmortem checks path's five required sections. With path
// empty, it resolves to the most recently created postmortem (design.md
// Decision 4); with a path given, it validates that file and promotes
// whichever attempt owns it, not necessarily the most recent one
// (design.md Decision 5). On success it promotes the owning attempt
// from "solved" to "complete" (a no-op if it's already "complete").
func (r *Runner) ValidatePostmortem(ctx context.Context, path string) (ValidationResult, error) {
	pf, err := progress.Load(r.ProgressPath)
	if err != nil {
		return ValidationResult{}, err
	}

	readPath := path
	if path == "" {
		index, ok := latestPostmortemAttempt(pf)
		if !ok {
			return ValidationResult{}, errors.New("no postmortem to validate; run `breaklab postmortem` first")
		}
		path = pf.Attempts[index].Postmortem
		// A path that came from progress.json must stay inside the
		// postmortems directory. An explicit CLI argument is the user's
		// own choice and is read as given.
		confined, err := postmortem.Confined(r.PostmortemsDir, path)
		if err != nil {
			return ValidationResult{}, err
		}
		readPath = confined
	}

	data, err := os.ReadFile(readPath)
	if err != nil {
		return ValidationResult{}, err
	}
	_, body, err := postmortem.SplitFrontMatter(string(data))
	if err != nil {
		return ValidationResult{}, err
	}

	result := ValidationResult{RootCauseWarning: r.rootCauseWarning(pf, path, body)}

	if failing := postmortem.ValidateSections(body); len(failing) > 0 {
		result.Failing = failing
		return result, nil
	}
	result.Valid = true

	for i, a := range pf.Attempts {
		if a.Postmortem == path {
			if err := progress.MarkComplete(r.ProgressPath, i); err != nil {
				return ValidationResult{}, err
			}
			break
		}
	}
	return result, nil
}

// rootCauseWarning returns an advisory message when body's Root cause
// section shares no meaningful word with the owning attempt's actual
// variant's hints, or "" when it does — or when the owning attempt or
// its variant can't be resolved, since this check must never block
// validation on a lookup failure (design.md - Goals).
func (r *Runner) rootCauseWarning(pf *progress.ProgressFile, path, body string) string {
	var attempt progress.Attempt
	found := false
	for _, a := range pf.Attempts {
		if a.Postmortem == path {
			attempt, found = a, true
			break
		}
	}
	if !found {
		return ""
	}

	lab, err := r.findLab(attempt.Lab)
	if err != nil {
		return ""
	}
	var hints []string
	for _, v := range lab.Variants {
		if v.ID == attempt.Variant {
			hints = v.Hints
			break
		}
	}
	if len(hints) == 0 {
		return ""
	}

	rootCause := postmortem.SplitIntoSections(body)["Root cause"]
	if postmortem.RootCauseMentionsVariant(rootCause, hints) {
		return ""
	}
	expected := postmortem.MeaningfulWords(strings.Join(hints, " "))
	return fmt.Sprintf("Root cause doesn't mention any of: %s — double check it describes what actually broke.", strings.Join(expected, ", "))
}

// newPostmortemPath returns postmortems/<date>-<lab>-<variant>.md (or
// postmortems/<date>-<lab>.md for a lab with no variant), adding a -2,
// -3, ... suffix if that name is already taken (design.md Decision 8).
func (r *Runner) newPostmortemPath(date, lab, variant string) (string, error) {
	base := fmt.Sprintf("%s-%s", date, lab)
	if variant != "" {
		base = fmt.Sprintf("%s-%s", base, variant)
	}
	for n := 1; ; n++ {
		name := base + ".md"
		if n > 1 {
			name = fmt.Sprintf("%s-%d.md", base, n)
		}
		path := filepath.Join(r.PostmortemsDir, name)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return path, nil
		} else if err != nil {
			return "", err
		}
	}
}
