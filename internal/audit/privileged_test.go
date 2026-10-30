package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// knownPrivilegedLabs is the deliberate, reviewed allowlist of labs that
// opt into --privileged (security spec: "Privileged containers are an
// explicit, per-lab opt-in"). Adding a lab here is a conscious decision:
// if this test fails because a lab.yaml sets privileged: true without
// being listed, add it here only after confirming why it actually needs
// the extra container access.
const k3sReason = "runs a real single-node k3s control plane, which needs cgroup access an unprivileged container is denied"

var knownPrivilegedLabs = map[string]string{
	"backend-dependency-down":       k3sReason,
	"cache-dependency-down":         k3sReason,
	"crash-loop-worker":             k3sReason,
	"deployment-never-ready":        k3sReason,
	"image-pull-backoff":            k3sReason,
	"impossible-scheduling-request": k3sReason,
	"k8s-blind-spot":                k3sReason,
	"k8sa-001":                      k3sReason,
	"k8sa-002":                      k3sReason,
	"k8sa-003":                      k3sReason,
	"k8sd-001":                      k3sReason,
	"k8sd-002":                      k3sReason,
	"k8sd-003":                      k3sReason,
	"k8sd-004":                      k3sReason,
	"k8sd-005":                      k3sReason,
	"k8sd-006":                      k3sReason,
	"k8sd-007":                      k3sReason,
	"k8sd-008":                      k3sReason,
	"k8sf-001":                      k3sReason,
	"k8sf-002":                      k3sReason,
	"k8sf-003":                      k3sReason,
	"k8sf-004":                      k3sReason,
	"k8sf-005":                      k3sReason,
	"k8sf-006":                      k3sReason,
	"k8ss-001":                      k3sReason,
	"k8ss-002":                      k3sReason,
	"k8ss-003":                      k3sReason,
	"k8ss-004":                      k3sReason,
	"network-policy-block":          k3sReason,
	"oom-killed":                    k3sReason,
	"pvc-bad-storageclass":          k3sReason,
	"ready-but-not-ready":           k3sReason,
	"selector-mismatch":             k3sReason,
	"stale-config":                  k3sReason,
	"stale-configmap":               k3sReason,
}

// TestPrivilegedLabsAreDeliberateAndDocumented enforces the security
// spec's "Privileged containers are an explicit, per-lab opt-in"
// requirement: every lab.yaml setting environment.privileged: true must
// be a known, reviewed entry with a documented reason, and every known
// entry must still correspond to a real lab.
func TestPrivilegedLabsAreDeliberateAndDocumented(t *testing.T) {
	labsDir := filepath.Join("..", "..", "labs")
	entries, err := os.ReadDir(labsDir)
	if err != nil {
		t.Fatalf("reading labs dir: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(labsDir, e.Name(), "lab.yaml"))
		if err != nil {
			continue
		}
		privileged := false
		for _, line := range strings.Split(string(data), "\n") {
			if strings.TrimSpace(line) == "privileged: true" {
				privileged = true
				break
			}
		}
		if !privileged {
			continue
		}
		reason, known := knownPrivilegedLabs[e.Name()]
		if !known || strings.TrimSpace(reason) == "" {
			t.Errorf("lab %q sets environment.privileged: true but isn't in knownPrivilegedLabs with a documented reason", e.Name())
		}
	}

	for name := range knownPrivilegedLabs {
		if _, err := os.Stat(filepath.Join(labsDir, name, "lab.yaml")); err != nil {
			t.Errorf("knownPrivilegedLabs lists %q but its lab.yaml doesn't exist", name)
		}
	}
}
