package webcore

import (
	"fmt"
	"path/filepath"
	"strings"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

// TicketID formats attemptIndex — its position within progress.json's
// Attempts slice — as a stable, human-readable display id shown
// consistently everywhere that attempt is referenced. It is a
// display-layer value only: it introduces no new stored field and is
// never used in place of a lab or variant id for any Docker,
// filesystem, or routing operation.
func TicketID(attemptIndex int) string {
	return fmt.Sprintf("TKT-%04d", 1000+attemptIndex)
}

// SlugFor derives a postmortem's URL slug from its file path. A trivial
// inline copy, not shared with internal/site (design.md Decision 3).
func SlugFor(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// PortNodeView is one of the active container's declared ports.
type PortNodeView struct {
	HostPort      string
	ContainerPort string
	OpenHref      string // best-effort "open" link; see proposal.md's Topology view bullet
}

// LabRowView is one lab's row, shared by the roadmap, challenge
// browser, and chaos views.
type LabRowView struct {
	ID               string
	Title            string
	Summary          string
	StageID          string
	StageTitle       string
	Status           string
	Difficulty       int
	EstimatedMinutes int
	Tags             []string
	TicketID         string // set only for a lab with at least one recorded attempt
	PostmortemHref   string // set only for a complete lab with a recorded postmortem
	ReviewDue        bool   // true for a complete lab whose last solve is past progress.ReviewInterval

	// IsUpNext marks the single lab, across the entire roadmap in
	// roadmap.yaml/lab.yaml order, that is the first one not yet
	// "complete" — the room-list's one "Up next" cursor. Set only by the
	// roadmap view; the challenge browser's rows aren't in roadmap
	// order, so it never sets this.
	IsUpNext bool
}

// HasTag reports whether tags contains tag.
func HasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// LabRow builds lab's row within group, given the catalog's already
// applied status (progress.ApplyStatus), the lab's latest complete
// attempt (if any), and which labs are currently review-due
// (progress.ReviewDue).
func LabRow(group catalog.StageGroup, lab catalog.Lab, latest map[string]progress.Attempt, pf *progress.ProgressFile, reviewDue map[string]bool) LabRowView {
	row := LabRowView{
		ID:               lab.ID,
		Title:            lab.Title,
		Summary:          lab.Summary,
		StageID:          group.Stage.ID,
		StageTitle:       group.Stage.Title,
		Status:           lab.Status,
		Difficulty:       lab.Difficulty,
		EstimatedMinutes: lab.EstimatedMinutes,
		Tags:             lab.Tags,
		ReviewDue:        reviewDue[lab.ID],
	}
	if attempt, ok := latest[lab.ID]; ok && attempt.Postmortem != "" {
		row.PostmortemHref = "/postmortems/" + SlugFor(attempt.Postmortem)
	}
	if idx, ok := progress.LatestAttemptIndex(pf, lab.ID); ok {
		row.TicketID = TicketID(idx)
	}
	return row
}
