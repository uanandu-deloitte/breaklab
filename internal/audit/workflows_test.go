package audit

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	usesLine   = regexp.MustCompile(`^\s*-?\s*uses:\s*(\S+)`)
	pinnedUses = regexp.MustCompile(`^[\w.-]+/[\w./-]+@[0-9a-f]{40}$`)
)

// Every action is pinned to a full commit SHA (tags are mutable), and no
// workflow runs untrusted pull-request code with pull_request_target.
func TestWorkflowsPinActionsAndAvoidPullRequestTarget(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(repoRoot, ".github", "workflows", "*.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no workflow files found")
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "pull_request_target") {
			t.Errorf("%s uses pull_request_target", f)
		}
		for i, line := range strings.Split(string(data), "\n") {
			m := usesLine.FindStringSubmatch(line)
			if m == nil || strings.HasPrefix(m[1], "./") {
				continue
			}
			if !pinnedUses.MatchString(m[1]) {
				t.Errorf("%s:%d: %q is not pinned to a 40-character commit SHA", filepath.Base(f), i+1, m[1])
			}
		}
	}
}
