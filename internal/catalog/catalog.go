package catalog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// ReadTicket reads labID's task brief: README.md for a lab using the
// curriculum doc's standalone-lab layout, falling back to the legacy
// TICKET.md every earlier lab still uses. Exactly one of the two is
// expected to exist; an unrelated read error (not "file doesn't exist")
// is returned immediately rather than silently falling through. Shared
// by internal/runner (Start prints the ticket to the CLI) and
// internal/web (the lab detail page), so the fallback stays in sync
// across both.
func ReadTicket(labsDir, labID string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(labsDir, labID, "README.md"))
	if err == nil {
		return data, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return os.ReadFile(filepath.Join(labsDir, labID, "TICKET.md"))
}

// LoadCatalog loads roadmap.yaml and every labs/<id>/lab.yaml under labsDir,
// validates each lab, and returns the valid labs grouped by roadmap stage in
// roadmap order, plus one warning per lab that failed to load or validate.
// A lab with a problem does not stop the rest of the catalog from loading.
func LoadCatalog(labsDir, roadmapPath string) (*Catalog, error) {
	roadmap, err := loadRoadmap(roadmapPath)
	if err != nil {
		return nil, err
	}

	stages := make(map[string]bool, len(roadmap.Stages))
	for _, s := range roadmap.Stages {
		stages[s.ID] = true
	}

	entries, err := os.ReadDir(labsDir)
	if err != nil {
		return nil, fmt.Errorf("reading labs directory %s: %w", labsDir, err)
	}

	var loaded []Lab
	var warnings []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		folder := entry.Name()

		lab, err := loadLab(filepath.Join(labsDir, folder))
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s: %v", folder, err))
			continue
		}

		if reason := validateLab(lab, folder, stages); reason != "" {
			warnings = append(warnings, fmt.Sprintf("%s: %s", folder, reason))
			continue
		}

		lab.Status = computeStatus(lab)
		loaded = append(loaded, lab)
	}

	// Prerequisites reference other labs by id, which can only be
	// checked once every lab.yaml has loaded — a lab naming an unknown
	// prerequisite is dropped the same non-blocking way any other
	// invalid lab is, rather than silently keeping a dangling reference
	// that could never unlock.
	validIDs := make(map[string]bool, len(loaded))
	for _, lab := range loaded {
		validIDs[lab.ID] = true
	}

	labsByStage := make(map[string][]Lab)
	for _, lab := range loaded {
		unknown := ""
		for _, id := range lab.Prerequisites {
			if !validIDs[id] {
				unknown = id
				break
			}
		}
		if unknown != "" {
			warnings = append(warnings, fmt.Sprintf("%s: unknown prerequisite %q", lab.ID, unknown))
			continue
		}
		labsByStage[lab.Stage] = append(labsByStage[lab.Stage], lab)
	}

	groups := make([]StageGroup, 0, len(roadmap.Stages))
	for _, stage := range roadmap.Stages {
		labs := labsByStage[stage.ID]
		// Easiest first within a stage — ID as a tiebreaker only for
		// labs tied on difficulty, so the order stays deterministic.
		sort.Slice(labs, func(i, j int) bool {
			if labs[i].Difficulty != labs[j].Difficulty {
				return labs[i].Difficulty < labs[j].Difficulty
			}
			return labs[i].ID < labs[j].ID
		})
		groups = append(groups, StageGroup{Stage: stage, Labs: labs})
	}

	return &Catalog{Stages: groups, Warnings: warnings}, nil
}

// computeStatus returns the status shown for a lab in `breaklab list`.
// Progress tracking (progress/progress.json) doesn't exist until a later
// change, so every lab is "available" until then.
func computeStatus(lab Lab) string {
	return "available"
}

func loadRoadmap(path string) (Roadmap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Roadmap{}, fmt.Errorf("reading roadmap.yaml: %w", err)
	}

	var rm Roadmap
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&rm); err != nil {
		return Roadmap{}, fmt.Errorf("parsing roadmap.yaml: %w", err)
	}

	capstones := 0
	for _, s := range rm.Stages {
		if !ValidID(s.ID) {
			return Roadmap{}, fmt.Errorf("roadmap.yaml: invalid stage id %q", s.ID)
		}
		if s.Capstone {
			capstones++
		}
	}
	if capstones > 1 {
		return Roadmap{}, fmt.Errorf("roadmap.yaml: only one stage may set capstone: true, found %d", capstones)
	}

	return rm, nil
}

func loadLab(dir string) (Lab, error) {
	path := filepath.Join(dir, "lab.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return Lab{}, fmt.Errorf("reading lab.yaml: %w", err)
	}

	var lab Lab
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&lab); err != nil {
		return Lab{}, fmt.Errorf("parsing lab.yaml: %w", err)
	}

	return lab, nil
}
