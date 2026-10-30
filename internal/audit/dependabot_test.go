package audit

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type dependabotConfig struct {
	Version int                `yaml:"version"`
	Updates []dependabotUpdate `yaml:"updates"`
}

type dependabotUpdate struct {
	PackageEcosystem string `yaml:"package-ecosystem"`
}

// Pinned GitHub Actions and the go.mod toolchain version only stay current
// if something proposes the bump; Dependabot has to still be configured
// for both ecosystems.
func TestDependabotCoversActionsAndGoToolchain(t *testing.T) {
	path := filepath.Join(repoRoot, ".github", "dependabot.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	var cfg dependabotConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("%s is not valid YAML: %v", path, err)
	}

	var sawActions, sawGomod bool
	for _, u := range cfg.Updates {
		switch u.PackageEcosystem {
		case "github-actions":
			sawActions = true
		case "gomod":
			sawGomod = true
		}
	}
	if !sawActions {
		t.Errorf("%s has no github-actions update entry", path)
	}
	if !sawGomod {
		t.Errorf("%s has no gomod update entry", path)
	}
}
