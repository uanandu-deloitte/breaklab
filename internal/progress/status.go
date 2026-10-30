package progress

import "breaklab/internal/catalog"

// Lab statuses computed by ApplyStatus.
const (
	LabStatusAvailable  = "available"
	LabStatusInProgress = "in-progress"
	LabStatusComplete   = "complete"

	// LabStatusLocked overrides whatever ApplyStatus's first pass
	// computed, for a lab whose Prerequisites aren't all "complete"
	// yet. A lab with no Prerequisites is never locked.
	LabStatusLocked = "locked"
)

// ApplyStatus computes each lab's status from pf and writes it onto
// cat's Lab.Status fields: "complete" if the lab has any "complete"
// attempt; "in-progress" if it has any "active" or "solved" attempt and
// no "complete" one; "available" otherwise — including a lab whose only
// attempts are all "abandoned" (nothing about it is still open).
//
// This lives here, not in internal/catalog, because internal/progress
// already imports internal/catalog (for StageCompletion); the reverse
// import would be a cycle.
func ApplyStatus(cat *catalog.Catalog, pf *ProgressFile) {
	complete := make(map[string]bool)
	inProgress := make(map[string]bool)
	for _, a := range pf.Attempts {
		switch a.Status {
		case StatusComplete, StatusSkipped:
			complete[a.Lab] = true
		case StatusActive, StatusSolved:
			inProgress[a.Lab] = true
		}
	}

	for i := range cat.Stages {
		for j := range cat.Stages[i].Labs {
			lab := &cat.Stages[i].Labs[j]
			switch {
			case complete[lab.ID]:
				lab.Status = LabStatusComplete
			case inProgress[lab.ID]:
				lab.Status = LabStatusInProgress
			default:
				lab.Status = LabStatusAvailable
			}
		}
	}

	applyPrerequisiteLocks(cat)
}

// applyPrerequisiteLocks overrides the status of every lab whose
// Prerequisites aren't all "complete" yet (per the first pass above) to
// LabStatusLocked. A lab with no Prerequisites is untouched — every lab
// shipped before Prerequisites existed has none, so this is a no-op for
// all of them.
func applyPrerequisiteLocks(cat *catalog.Catalog) {
	statusByID := make(map[string]string)
	for _, group := range cat.Stages {
		for _, lab := range group.Labs {
			statusByID[lab.ID] = lab.Status
		}
	}

	for i := range cat.Stages {
		for j := range cat.Stages[i].Labs {
			lab := &cat.Stages[i].Labs[j]
			if len(lab.Prerequisites) == 0 || lab.Status == LabStatusComplete {
				continue
			}
			for _, id := range lab.Prerequisites {
				if statusByID[id] != LabStatusComplete {
					lab.Status = LabStatusLocked
					break
				}
			}
		}
	}
}
