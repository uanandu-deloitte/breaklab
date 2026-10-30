package audit

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const repoRoot = "../.."

// Session state and the generated site must never be committed.
func TestGitignoreCoversSessionAndSite(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(repoRoot, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	entries := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		entries[strings.TrimSpace(line)] = true
	}
	// Anchored ("/site/"), not bare "site/": a bare pattern matches at any
	// depth and would collaterally hide the real, committed source trees
	// templates/site/ and internal/site/, not just the generated output
	// directory this rule exists for.
	for _, want := range []string{".breaklab/", "/site/"} {
		if !entries[want] {
			t.Errorf(".gitignore is missing %q", want)
		}
	}

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	out, err := exec.Command("git", "-C", repoRoot, "ls-files", "--", ".breaklab", "site").Output()
	if err != nil {
		t.Skipf("not inside a git work tree: %v", err)
	}
	if tracked := strings.TrimSpace(string(out)); tracked != "" {
		t.Errorf("files under .breaklab/ or the root site/ are tracked by git:\n%s", tracked)
	}

	// The source trees must stay trackable even though their path
	// contains "site" — this is what the bare "site/" pattern broke.
	for _, path := range []string{"templates/site", "internal/site"} {
		out, err := exec.Command("git", "-C", repoRoot, "check-ignore", path).Output()
		if err == nil && strings.TrimSpace(string(out)) != "" {
			t.Errorf("%s is ignored by .gitignore but is real source, not generated output", path)
		}
	}
}
