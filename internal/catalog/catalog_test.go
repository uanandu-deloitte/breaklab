package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testRoadmap = `
stages:
  - id: linux
    title: Linux fundamentals
  - id: git
    title: Git
`

const validLabYAML = `
id: web-down
title: The website returns nothing
stage: linux
difficulty: 1
estimated_minutes: 20
environment:
  dockerfile: Dockerfile
  ports: ["8080:80"]
check:
  script: check.sh
  run_in: host
hints:
  - Start with the symptom.
variants:
  - id: config-typo
    script: variants/config-typo.sh
    weight: 1
    hints:
      - Ask nginx to test its own configuration.
`

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// Task 3.1: roadmap.yaml and lab.yaml decode without error under strict decoding.
func TestLoadRoadmapAndLab(t *testing.T) {
	dir := t.TempDir()
	roadmapPath := filepath.Join(dir, "roadmap.yaml")
	writeFile(t, roadmapPath, testRoadmap)

	labsDir := filepath.Join(dir, "labs")
	writeFile(t, filepath.Join(labsDir, "web-down", "lab.yaml"), validLabYAML)

	cat, err := LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if len(cat.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", cat.Warnings)
	}

	var found bool
	for _, group := range cat.Stages {
		if group.Stage.ID != "linux" {
			continue
		}
		for _, lab := range group.Labs {
			if lab.ID == "web-down" {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("expected web-down under linux stage, got %+v", cat.Stages)
	}
}

// Task 3.2: id, stage, and script-path validation.
func TestValidateLab(t *testing.T) {
	stages := map[string]bool{"linux": true}

	base := Lab{
		ID:    "web-down",
		Stage: "linux",
		Check: Check{Script: "check.sh"},
		Variants: []Variant{
			{ID: "config-typo", Script: "variants/config-typo.sh"},
		},
	}

	if reason := validateLab(base, "web-down", stages); reason != "" {
		t.Fatalf("expected valid lab to pass, got reason: %s", reason)
	}

	t.Run("bad id", func(t *testing.T) {
		lab := base
		lab.ID = "Web_Down"
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected invalid id to fail validation")
		}
	})

	t.Run("bad stage reference", func(t *testing.T) {
		lab := base
		lab.Stage = "does-not-exist"
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected unknown stage to fail validation")
		}
	})

	t.Run("escaping script path", func(t *testing.T) {
		lab := base
		lab.Check.Script = "../../etc/passwd"
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected escaping script path to fail validation")
		}
	})

	t.Run("bad tag", func(t *testing.T) {
		lab := base
		lab.Tags = []string{"Web Servers"}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected invalid tag to fail validation")
		}
	})

	t.Run("valid tags", func(t *testing.T) {
		lab := base
		lab.Tags = []string{"linux", "web-servers"}
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected valid tags to pass, got reason: %s", reason)
		}
	})

	t.Run("valid volume", func(t *testing.T) {
		lab := base
		lab.Environment.Volumes = []Volume{{Name: "data", ContainerPath: "/data"}}
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected valid volume to pass, got reason: %s", reason)
		}
	})

	t.Run("bad volume name", func(t *testing.T) {
		lab := base
		lab.Environment.Volumes = []Volume{{Name: "My Data", ContainerPath: "/data"}}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected invalid volume name to fail validation")
		}
	})

	t.Run("duplicate volume name", func(t *testing.T) {
		lab := base
		lab.Environment.Volumes = []Volume{
			{Name: "data", ContainerPath: "/data"},
			{Name: "data", ContainerPath: "/other"},
		}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected duplicate volume name to fail validation")
		}
	})

	t.Run("relative container path", func(t *testing.T) {
		lab := base
		lab.Environment.Volumes = []Volume{{Name: "data", ContainerPath: "data"}}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected relative container_path to fail validation")
		}
	})

	t.Run("container path with dot-dot still resolves safely", func(t *testing.T) {
		// filepath.Clean fully resolves ".." for an absolute path (it
		// can't go above root), so this normalizes to the harmless
		// "/etc" rather than escaping anywhere.
		lab := base
		lab.Environment.Volumes = []Volume{{Name: "data", ContainerPath: "/data/../../etc"}}
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected this to normalize to a valid path, got reason: %s", reason)
		}
	})

	t.Run("container path is root", func(t *testing.T) {
		lab := base
		lab.Environment.Volumes = []Volume{{Name: "data", ContainerPath: "/"}}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected root container_path to fail validation")
		}
	})

	t.Run("negative pids_limit", func(t *testing.T) {
		lab := base
		lab.Environment.PidsLimit = -1
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected negative pids_limit to fail validation")
		}
	})

	t.Run("positive pids_limit", func(t *testing.T) {
		lab := base
		lab.Environment.PidsLimit = 32
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected positive pids_limit to pass, got reason: %s", reason)
		}
	})

	t.Run("valid loop", func(t *testing.T) {
		lab := base
		lab.Loop = Loop{
			Learn:  "A primer.",
			Guided: []GuidedStep{{Do: "curl -I localhost", Why: "See the symptom first."}},
			Practice: PracticeBlock{
				Task:      "Rehearse it.",
				Questions: []PracticeQuestion{{Prompt: "Which command?", Answer: "curl -I"}},
			},
		}
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected valid loop to pass, got reason: %s", reason)
		}
	})

	t.Run("guided step missing why", func(t *testing.T) {
		lab := base
		lab.Loop.Guided = []GuidedStep{{Do: "curl -I localhost", Why: ""}}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected guided step missing why to fail validation")
		}
	})

	t.Run("practice question missing answer", func(t *testing.T) {
		lab := base
		lab.Loop.Practice.Questions = []PracticeQuestion{{Prompt: "Which command?", Answer: ""}}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected practice question missing answer to fail validation")
		}
	})

	t.Run("plural checks without legacy check is valid", func(t *testing.T) {
		lab := base
		lab.Check = Check{}
		lab.Checks = []NamedCheck{{ID: "app-can-read", Script: "checks/app-can-read.sh"}}
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected plural checks to satisfy validation, got reason: %s", reason)
		}
	})

	t.Run("neither check nor checks set", func(t *testing.T) {
		lab := base
		lab.Check = Check{}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected a lab with neither check nor checks to fail validation")
		}
	})

	t.Run("duplicate check id", func(t *testing.T) {
		lab := base
		lab.Check = Check{}
		lab.Checks = []NamedCheck{
			{ID: "app-can-read", Script: "checks/a.sh"},
			{ID: "app-can-read", Script: "checks/b.sh"},
		}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected duplicate check id to fail validation")
		}
	})

	t.Run("invalid check script path", func(t *testing.T) {
		lab := base
		lab.Check = Check{}
		lab.Checks = []NamedCheck{{ID: "app-can-read", Script: "../escape.sh"}}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected escaping checks script path to fail validation")
		}
	})

	t.Run("valid type", func(t *testing.T) {
		lab := base
		lab.Type = "guided"
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected valid type to pass, got reason: %s", reason)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		lab := base
		lab.Type = "tutorial"
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected invalid type to fail validation")
		}
	})

	t.Run("valid lesson path", func(t *testing.T) {
		lab := base
		lab.Lesson = "lesson.md"
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected valid lesson path to pass, got reason: %s", reason)
		}
	})

	t.Run("escaping lesson path", func(t *testing.T) {
		lab := base
		lab.Lesson = "../lesson.md"
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected escaping lesson path to fail validation")
		}
	})

	t.Run("valid prerequisite id", func(t *testing.T) {
		lab := base
		lab.Prerequisites = []string{"lnx-001"}
		if reason := validateLab(lab, "web-down", stages); reason != "" {
			t.Fatalf("expected valid prerequisite id to pass, got reason: %s", reason)
		}
	})

	t.Run("invalid prerequisite id", func(t *testing.T) {
		lab := base
		lab.Prerequisites = []string{"LNX_001"}
		if reason := validateLab(lab, "web-down", stages); reason == "" {
			t.Fatal("expected invalid prerequisite id to fail validation")
		}
	})
}

// Task 3.3: one invalid lab produces a warning without blocking the rest.
func TestLoadCatalogInvalidLabScenario(t *testing.T) {
	dir := t.TempDir()
	roadmapPath := filepath.Join(dir, "roadmap.yaml")
	writeFile(t, roadmapPath, testRoadmap)

	labsDir := filepath.Join(dir, "labs")
	writeFile(t, filepath.Join(labsDir, "web-down", "lab.yaml"), validLabYAML)

	brokenLabYAML := `
id: broken-lab
title: A lab with a missing stage
stage: does-not-exist
difficulty: 1
estimated_minutes: 10
environment:
  dockerfile: Dockerfile
  ports: ["8081:80"]
check:
  script: check.sh
  run_in: host
`
	writeFile(t, filepath.Join(labsDir, "broken-lab", "lab.yaml"), brokenLabYAML)

	cat, err := LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if len(cat.Warnings) != 1 {
		t.Fatalf("expected exactly one warning, got %v", cat.Warnings)
	}

	var found bool
	for _, group := range cat.Stages {
		for _, lab := range group.Labs {
			if lab.ID == "web-down" {
				found = true
			}
			if lab.ID == "broken-lab" {
				t.Fatal("broken-lab should not have loaded")
			}
		}
	}
	if !found {
		t.Fatal("expected web-down to still load alongside the broken lab")
	}
}

// TestLoadCatalogDropsLabWithUnknownPrerequisite confirms a lab naming a
// prerequisite id that doesn't exist anywhere in the catalog is dropped
// with a warning, the same non-blocking way any other invalid lab is,
// while a lab whose prerequisite does exist loads normally.
func TestLoadCatalogDropsLabWithUnknownPrerequisite(t *testing.T) {
	dir := t.TempDir()
	roadmapPath := filepath.Join(dir, "roadmap.yaml")
	writeFile(t, roadmapPath, testRoadmap)

	labsDir := filepath.Join(dir, "labs")
	writeFile(t, filepath.Join(labsDir, "web-down", "lab.yaml"), validLabYAML)
	writeFile(t, filepath.Join(labsDir, "dangling-prereq", "lab.yaml"), `
id: dangling-prereq
title: Depends on nothing real
stage: linux
difficulty: 1
estimated_minutes: 10
environment:
  dockerfile: Dockerfile
  ports: ["8082:80"]
check:
  script: check.sh
  run_in: host
prerequisites:
  - no-such-lab
`)
	writeFile(t, filepath.Join(labsDir, "real-prereq", "lab.yaml"), `
id: real-prereq
title: Depends on web-down
stage: linux
difficulty: 1
estimated_minutes: 10
environment:
  dockerfile: Dockerfile
  ports: ["8083:80"]
check:
  script: check.sh
  run_in: host
prerequisites:
  - web-down
`)

	cat, err := LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if len(cat.Warnings) != 1 || !strings.Contains(cat.Warnings[0], "unknown prerequisite") {
		t.Fatalf("expected exactly one unknown-prerequisite warning, got %v", cat.Warnings)
	}

	var sawDangling, sawReal bool
	for _, group := range cat.Stages {
		for _, lab := range group.Labs {
			if lab.ID == "dangling-prereq" {
				sawDangling = true
			}
			if lab.ID == "real-prereq" {
				sawReal = true
			}
		}
	}
	if sawDangling {
		t.Error("dangling-prereq should not have loaded")
	}
	if !sawReal {
		t.Error("real-prereq should have loaded: its prerequisite exists")
	}
}

// add-random-variants task 4.2: web-down's declared variants (and check
// script) must each resolve to a real file on disk, not just a
// well-formed path — change 1's validation only checked the shape of the
// path string. Uses the repo's real labs/ and roadmap.yaml rather than a
// fixture, since the point is to catch a missing script for the lab as
// shipped.
func TestWebDownVariantScriptsExistOnDisk(t *testing.T) {
	labsDir := filepath.Join("..", "..", "labs")
	roadmapPath := filepath.Join("..", "..", "roadmap.yaml")

	cat, err := LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if len(cat.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", cat.Warnings)
	}

	var lab Lab
	var found bool
	for _, group := range cat.Stages {
		for _, l := range group.Labs {
			if l.ID == "web-down" {
				lab, found = l, true
			}
		}
	}
	if !found {
		t.Fatal("web-down did not load")
	}
	if len(lab.Variants) < 2 {
		t.Fatalf("expected web-down to declare 2+ variants, got %d", len(lab.Variants))
	}

	labDir := filepath.Join(labsDir, lab.ID)
	assertExists := func(rel string) {
		path := filepath.Join(labDir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("declared script %s does not exist on disk: %v", path, err)
		}
	}

	assertExists(lab.Check.Script)
	for _, v := range lab.Variants {
		assertExists(v.Script)
		assertExists(filepath.Join("solutions", v.ID+".sh"))
	}
}

// Task 3.4: status computation defaults to "available" with no progress data.
func TestComputeStatusAvailable(t *testing.T) {
	dir := t.TempDir()
	roadmapPath := filepath.Join(dir, "roadmap.yaml")
	writeFile(t, roadmapPath, testRoadmap)

	labsDir := filepath.Join(dir, "labs")
	writeFile(t, filepath.Join(labsDir, "web-down", "lab.yaml"), validLabYAML)

	cat, err := LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}

	var status string
	for _, group := range cat.Stages {
		for _, lab := range group.Labs {
			if lab.ID == "web-down" {
				status = lab.Status
			}
		}
	}
	if status != "available" {
		t.Fatalf("expected status %q, got %q", "available", status)
	}
}

// Tags are optional metadata used by the challenge browser's topic filter;
// LoadCatalog must decode and preserve them.
func TestLoadCatalogParsesTags(t *testing.T) {
	dir := t.TempDir()
	roadmapPath := filepath.Join(dir, "roadmap.yaml")
	writeFile(t, roadmapPath, testRoadmap)

	labsDir := filepath.Join(dir, "labs")
	writeFile(t, filepath.Join(labsDir, "web-down", "lab.yaml"), validLabYAML+"tags: [linux, web-servers]\n")

	cat, err := LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if len(cat.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", cat.Warnings)
	}

	var tags []string
	for _, group := range cat.Stages {
		for _, lab := range group.Labs {
			if lab.ID == "web-down" {
				tags = lab.Tags
			}
		}
	}
	want := []string{"linux", "web-servers"}
	if len(tags) != len(want) || tags[0] != want[0] || tags[1] != want[1] {
		t.Fatalf("tags = %v, want %v", tags, want)
	}
}

// An invalid tag makes the whole lab invalid, same as a bad id or path.
func TestLoadCatalogInvalidTagScenario(t *testing.T) {
	dir := t.TempDir()
	roadmapPath := filepath.Join(dir, "roadmap.yaml")
	writeFile(t, roadmapPath, testRoadmap)

	labsDir := filepath.Join(dir, "labs")
	writeFile(t, filepath.Join(labsDir, "web-down", "lab.yaml"), validLabYAML+"tags: [Web_Servers]\n")

	cat, err := LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if len(cat.Warnings) != 1 {
		t.Fatalf("expected exactly one warning, got %v", cat.Warnings)
	}
	for _, group := range cat.Stages {
		for _, lab := range group.Labs {
			if lab.ID == "web-down" {
				t.Fatal("web-down with an invalid tag should not have loaded")
			}
		}
	}
}

func TestValidID(t *testing.T) {
	for _, ok := range []string{"web-down", "config-typo", "linux", "a1"} {
		if !ValidID(ok) {
			t.Errorf("ValidID(%q) = false, want true", ok)
		}
	}
	for _, bad := range []string{"", "..", "a/b", "Web_Down", "../etc", "a b", "x.md"} {
		if ValidID(bad) {
			t.Errorf("ValidID(%q) = true, want false", bad)
		}
	}
}

func TestReadTicketPrefersReadmeOverLegacyTicket(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lnx-001", "README.md"), "from readme")
	data, err := ReadTicket(dir, "lnx-001")
	if err != nil || string(data) != "from readme" {
		t.Fatalf("data = %q, err = %v", data, err)
	}
}

func TestReadTicketFallsBackToLegacyTicket(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "web-down", "TICKET.md"), "from ticket")
	data, err := ReadTicket(dir, "web-down")
	if err != nil || string(data) != "from ticket" {
		t.Fatalf("data = %q, err = %v", data, err)
	}
}

func TestReadTicketErrorsWhenNeitherExists(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "empty-lab"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadTicket(dir, "empty-lab"); err == nil {
		t.Fatal("expected an error when neither README.md nor TICKET.md exists")
	}
}
