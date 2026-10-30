package catalog

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var idPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

// ValidID reports whether s is a safe lab, variant, or stage id: one or
// more of [a-z0-9-], so it can never carry a path separator or "..".
func ValidID(s string) bool {
	return idPattern.MatchString(s)
}

// validScriptPath reports whether p is a non-empty relative path that stays
// inside the lab's own directory: no absolute paths and no ".." segments.
func validScriptPath(p string) bool {
	if p == "" || filepath.IsAbs(p) {
		return false
	}
	clean := filepath.ToSlash(filepath.Clean(p))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return false
	}
	return true
}

// validContainerPath reports whether p is a safe absolute in-container
// mount point and not the root itself. Unlike validScriptPath, ".."
// segments need no separate check here: filepath.Clean always fully
// resolves them for an absolute path (it can't go above root), so
// there's nothing left to escape.
func validContainerPath(p string) bool {
	if p == "" || !filepath.IsAbs(p) {
		return false
	}
	return filepath.ToSlash(filepath.Clean(p)) != "/"
}

// validateLab checks a lab against the Lab catalog requirement: its id,
// stage reference, and script paths must all be valid. It returns an empty
// string when the lab is valid, or a one-line reason otherwise.
func validateLab(lab Lab, folder string, stages map[string]bool) string {
	if !ValidID(lab.ID) {
		return fmt.Sprintf("invalid lab id %q", lab.ID)
	}
	if lab.ID != folder {
		return fmt.Sprintf("lab id %q does not match folder name %q", lab.ID, folder)
	}
	if !ValidID(lab.Stage) || !stages[lab.Stage] {
		return fmt.Sprintf("stage %q not found in roadmap.yaml", lab.Stage)
	}
	// A lab verifies itself through either the legacy singular Check or
	// the newer plural Checks, never neither — but not both are
	// required to be set at once, since a lab using Checks leaves Check
	// at its zero value.
	if len(lab.Checks) == 0 {
		if !validScriptPath(lab.Check.Script) {
			return fmt.Sprintf("check script path %q is invalid", lab.Check.Script)
		}
	} else {
		seenCheckIDs := make(map[string]bool, len(lab.Checks))
		for i, c := range lab.Checks {
			if !ValidID(c.ID) {
				return fmt.Sprintf("checks[%d] has invalid id %q", i, c.ID)
			}
			if seenCheckIDs[c.ID] {
				return fmt.Sprintf("duplicate check id %q", c.ID)
			}
			seenCheckIDs[c.ID] = true
			if !validScriptPath(c.Script) {
				return fmt.Sprintf("checks[%d] %q script path %q is invalid", i, c.ID, c.Script)
			}
		}
	}

	seenVolumes := make(map[string]bool, len(lab.Environment.Volumes))
	for _, v := range lab.Environment.Volumes {
		if !ValidID(v.Name) {
			return fmt.Sprintf("invalid volume name %q", v.Name)
		}
		if seenVolumes[v.Name] {
			return fmt.Sprintf("duplicate volume name %q", v.Name)
		}
		seenVolumes[v.Name] = true
		if !validContainerPath(v.ContainerPath) {
			return fmt.Sprintf("volume %q container_path %q is invalid", v.Name, v.ContainerPath)
		}
	}
	if lab.Environment.PidsLimit < 0 {
		return fmt.Sprintf("pids_limit %d must not be negative", lab.Environment.PidsLimit)
	}

	seenVariants := make(map[string]bool, len(lab.Variants))
	for _, v := range lab.Variants {
		if !ValidID(v.ID) {
			return fmt.Sprintf("invalid variant id %q", v.ID)
		}
		if seenVariants[v.ID] {
			return fmt.Sprintf("duplicate variant id %q", v.ID)
		}
		seenVariants[v.ID] = true
		if !validScriptPath(v.Script) {
			return fmt.Sprintf("variant %q script path %q is invalid", v.ID, v.Script)
		}
	}

	for _, tag := range lab.Tags {
		if !ValidID(tag) {
			return fmt.Sprintf("invalid tag %q", tag)
		}
	}

	for i, step := range lab.Loop.Guided {
		if strings.TrimSpace(step.Do) == "" || strings.TrimSpace(step.Why) == "" {
			return fmt.Sprintf("loop.guided[%d] must set both do and why", i)
		}
	}
	for i, q := range lab.Loop.Practice.Questions {
		if strings.TrimSpace(q.Prompt) == "" || strings.TrimSpace(q.Answer) == "" {
			return fmt.Sprintf("loop.practice.questions[%d] must set both prompt and answer", i)
		}
	}

	if lab.Type != "" && !validLabTypes[lab.Type] {
		return fmt.Sprintf("invalid type %q", lab.Type)
	}
	if lab.Lesson != "" && !validScriptPath(lab.Lesson) {
		return fmt.Sprintf("lesson path %q is invalid", lab.Lesson)
	}
	for _, id := range lab.Prerequisites {
		if !ValidID(id) {
			return fmt.Sprintf("invalid prerequisite id %q", id)
		}
	}

	return ""
}

// validLabTypes are the curriculum doc's five standalone lab types
// (Lab Types section): guided, practice, break-fix, incident, project.
var validLabTypes = map[string]bool{
	"guided":    true,
	"practice":  true,
	"break-fix": true,
	"incident":  true,
	"project":   true,
}
