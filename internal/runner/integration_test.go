package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/docker"
	"breaklab/internal/progress"
)

// endToEndLabs is every lab this suite drives against real Docker, with
// its full variant list. Every lab in labs/ should have an entry here —
// TestAllLabsCovered fails if one is missing, so a newly authored lab
// can't silently ship without an end-to-end check.
var endToEndLabs = map[string][]string{
	"alert-never-fires":       {"wrong-operator", "wrong-threshold"},
	"backend-dependency-down": {"scaled-to-zero", "crash-looping"}, // real k3s control plane, privileged + slow
	"bind-mount-shadow":       {"root-owned-mount", "wrong-uid-mount"},
	"bisect-the-break":        {"strictness-exit", "bad-threshold"},
	"blind-after-rotation":    {"stale-path", "malformed-config"},
	"bloated-data":            {"runaway-log", "stale-temp-files"},
	"cache-dependency-down":   {"scaled-to-zero", "crash-looping"}, // real k3s control plane, privileged + slow
	"clock-skew-timeline":     {"clock-drifted", "ntp-misconfigured"},
	"cron-blindspot":          {"bare-command", "relative-path"},
	"crash-loop-worker":       {"bad-command", "missing-env"}, // real k3s control plane, privileged + slow
	"deleted-file-still-open": {"rotated-without-truncate", "stray-debug-process"},
	"deployment-never-ready":  {"bad-selector", "bad-readiness-probe"}, // real k3s control plane, privileged + slow
	"drifted-config":          {"stale-apply", "manual-drift"},
	"fake-green":              {"bad-releases-dir", "bad-activation-perms"},
	"git-lost-work":           {"hard-reset", "detached-head"},
	"half-shipped":            {"wrong-pin-hard-reset", "missing-module-detached-commit"},
	"k8s-blind-spot":          {"selector-drift-stale-metric", "bad-probe-alert-never-fires"}, // real k3s control plane, privileged + slow

	// git-00X labs are the version-control path's standalone-lab model
	// (same model as lnx-001/.../015 above): no variants/solutions
	// chaos-injection, driven by their own dedicated TestGit*EndToEnd
	// functions below instead of runLabEndToEnd — these empty slices
	// exist only to satisfy TestAllLabsCovered's "every lab has an
	// entry" check.
	"git-001":                       {},
	"git-002":                       {},
	"git-003":                       {},
	"git-004":                       {},
	"git-005":                       {},
	"git-006":                       {},
	"git-007":                       {},
	"git-008":                       {},
	"image-pull-backoff":            {"bad-tag", "unreachable-registry"},     // real k3s control plane, privileged + slow
	"impossible-scheduling-request": {"impossible-cpu", "impossible-memory"}, // real k3s control plane, privileged + slow
	"infra-env-override":            {"stray-bashrc-export", "wrapper-hardcodes-it"},
	"infra-stuck-lock":              {"bad-lock-dir-perms", "stale-lock"},
	"k8sa-001":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sa-002":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sa-003":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sd-001":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sd-002":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sd-003":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sd-004":                      {"scaled-to-zero", "crash-looping"}, // real k3s control plane, privileged + slow
	"k8sd-005":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sd-006":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sd-007":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sd-008":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sf-001":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sf-002":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sf-003":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sf-004":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sf-005":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8sf-006":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8ss-001":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8ss-002":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8ss-003":                      {"setup"},                           // real k3s control plane, privileged + slow
	"k8ss-004":                      {"setup"},                           // real k3s control plane, privileged + slow
	"leaked-secret":                 {"env-file", "inline-code"},
	"metric-rename-gap":             {"broken-alias-bridge", "dashboard-not-updated"},
	"net-001":                       {},
	"net-002":                       {},
	"net-003":                       {},
	"net-004":                       {},
	"net-005":                       {},
	"net-006":                       {},
	"net-007":                       {},
	"network-policy-block":          {"deny-all-ingress", "overbroad-selector"}, // real k3s control plane, privileged + slow
	"oom-killed":                    {"runaway-allocation", "limit-too-low"},    // real k3s control plane, privileged + slow, pulls python:3.12-alpine
	"path-shadow":                   {"fake-jq-stub", "legacy-tool-symlink"},
	"phantom-release":               {"merged-wrong-branch-stale-cache", "unpushed-tag-points-old-commit"},
	"pids-ceiling":                  {"no-concurrency-cap", "metrics-collector-leak"},
	"pinned-module-regression":      {"pinned-to-buggy", "pinned-to-missing"},
	"proxy-maze":                    {"bad-upstream", "bad-syntax"},
	"pvc-bad-storageclass":          {"missing-storageclass", "unsupported-access-mode"}, // real k3s control plane, privileged + slow
	"queue-drain":                   {"missing-config", "bad-permissions"},
	"quiet-shipper":                 {"missing-jq", "wrong-destination"},
	"ready-but-not-ready":           {"wrong-probe-path", "wrong-probe-port"}, // real k3s control plane, privileged + slow
	"selector-mismatch":             {"wrong-selector", "wrong-target-port"},  // real k3s control plane, privileged + slow
	"shipper-buffer-drop":           {"consumer-stopped", "queue-readonly"},
	"silent-pipeline":               {"hook-not-executable", "hook-broken-shebang"},
	"silent-stdout":                 {"log-to-file", "app-config-log-target"},
	"stale-cache-restore":           {"cache-never-invalidated", "cache-key-ignores-lockfile"},
	"stale-config":                  {"env-injected", "volume-cached"}, // real k3s control plane, privileged + slow
	"stale-configmap":               {"env-var", "subpath-mount"},      // real k3s control plane, privileged + slow
	"stale-stash":                   {"forgotten-stash", "wrong-stash-popped"},
	"stale-symlink":                 {"pointing-old", "dangling-link"},
	"stuck-lock":                    {"stale-lock", "bad-dir-perms"},
	"swallowed-signal":              {"piped-through-tee", "wrapper-script-no-exec"},
	"umask-mismatch":                {"service-env-umask", "shared-profile-umask"},
	"unexpanded-template":           {"bad-template", "bad-apply-script"},
	"unmanaged-resource":            {"never-imported", "path-mismatch"},
	"web-down":                      {"config-typo", "port-taken", "bad-perms"},
	"wrong-artifact-promoted":       {"config-names-wrong-artifact", "glob-latest-wins"},
	"zombie-pileup":                 {"leaked-subprocess", "missing-wait"},

	// obs-00X are the observability stage's foundational standalone-lab
	// model labs (same model as lnx-001/.../015 below): no
	// variants/solutions chaos-injection, driven by their own dedicated
	// TestObs*EndToEnd functions below instead of runLabEndToEnd.
	"obs-001": {},
	"obs-002": {},
	"obs-003": {},
	"obs-004": {},
	"obs-005": {},
	"obs-006": {},
	"obs-007": {},

	// lnx-001/003/005/010/015 are the curriculum doc's standalone-lab
	// model (Release 0.1): no variants/solutions chaos-injection, so
	// they're driven by their own dedicated TestLnx*EndToEnd functions
	// below instead of runLabEndToEnd — these empty slices exist only
	// to satisfy TestAllLabsCovered's "every lab has an entry" check.
	"lnx-001": {},
	"lnx-002": {},
	"lnx-003": {},
	"lnx-004": {},
	"lnx-005": {},
	"lnx-006": {},
	"lnx-007": {},
	"lnx-008": {},
	"lnx-009": {},
	"lnx-010": {},
	"lnx-011": {},
	"lnx-012": {},
	"lnx-013": {},
	"lnx-014": {},
	"lnx-015": {},
	"lnx-016": {},
	"lnx-017": {},
	"lnx-018": {},
	"lnx-019": {},
	"lnx-020": {},
	"lnx-021": {},
	"lnx-022": {},
	"lnx-023": {},
	"lnx-024": {},
	"lnx-025": {},

	// lnx-900 is the Linux stage's final incident: the older
	// multi-variant chaos-injection pattern (like web-down above), with
	// five seeded fault scenarios the runner picks between at random.
	"lnx-900": {"backend-not-running", "invalid-nginx-upstream", "config-permission-error", "port-collision", "disk-full"},

	// git-00X are the version-control path's standalone-lab-model labs
	// (same model as lnx-001/.../015 above): no variants/solutions
	// chaos-injection, driven by their own dedicated TestGit*EndToEnd
	// functions below instead of runLabEndToEnd.
	"git-009": {},
	"git-010": {},
	"git-011": {},
	"git-012": {},

	// web-001/web-002/web-003 and net-005/net-006/net-007 are the
	// networking stage's standalone-lab-model labs (same model as
	// lnx-001/.../015 above): no variants/solutions chaos-injection,
	// driven by their own dedicated TestWeb*EndToEnd/TestNet*EndToEnd
	// functions below instead of runLabEndToEnd.
	"web-001": {},
	"web-002": {},
	"web-003": {},

	// net-900 is the networking stage's final incident: the older
	// multi-variant chaos-injection pattern, same shape as web-down and
	// lnx-900 above, with five seeded fault scenarios across the
	// client -> nginx -> api -> database chain.
	"net-900": {"dns-misconfigured", "listener-down", "bad-upstream", "database-unreachable", "application-config-error"},

	// docker-00X are the Docker stage's foundational standalone-lab-model
	// labs (same model as lnx-001/.../015 above): no variants/solutions
	// chaos-injection, driven by their own dedicated TestDocker*EndToEnd
	// functions below instead of runLabEndToEnd.
	"docker-001": {},
	"docker-002": {},
	"docker-003": {},
	"docker-004": {},
	"docker-005": {},
	"docker-006": {},
	"docker-007": {},
	"docker-008": {},
	"docker-009": {},

	// cicd-00X are the CI/CD stage's foundational standalone-lab-model
	// labs (same model as lnx-001/.../015 above): no variants/solutions
	// chaos-injection, driven by their own dedicated TestCicd*EndToEnd
	// functions below instead of runLabEndToEnd.
	"cicd-001": {},
	"cicd-002": {},
	"cicd-003": {},
	"cicd-004": {},
	"cicd-005": {},
	"cicd-006": {},
	"cicd-007": {},

	// cicd-008/009 are real-tool standalone-lab-model labs, same
	// no-variants/solutions shape as cicd-00X above, but driving a real
	// `act` / `gitlab-ci-local` run instead of the toy pipeline runner
	// -- their own dedicated TestCicd*EndToEnd functions below.
	"cicd-008": {},
	"cicd-009": {},

	// iac-00X are the Infrastructure as code stage's foundational
	// standalone-lab-model labs (same model as lnx-001/.../015 above):
	// no variants/solutions chaos-injection, driven by their own
	// dedicated TestIac*EndToEnd functions below instead of
	// runLabEndToEnd.
	"iac-001": {},
	"iac-002": {},
	"iac-003": {},
	"iac-004": {},
	"iac-005": {},
	"iac-006": {},
	"iac-007": {},
	"iac-008": {},

	// iac-009/010/011 are real-tool standalone-lab-model labs, same
	// no-variants/solutions shape as iac-00X above, but driving real
	// `terraform` / `ansible-playbook` instead of the toy infra-apply
	// tool -- their own dedicated TestIac*EndToEnd functions below.
	"iac-009": {},
	"iac-010": {},
	"iac-011": {},

	// obs-008/009 are real-tool standalone-lab-model labs, same
	// no-variants/solutions shape as obs-00X below, but driving real
	// Prometheus/promtool and Grafana instead of the toy metric/alert
	// evaluator -- their own dedicated TestObs*EndToEnd functions below.
	"obs-008": {},
	"obs-009": {},

	// azure-00X are the new cloud stage's foundational standalone-lab-
	// model labs (same model as lnx-001/.../015 above): no
	// variants/solutions chaos-injection, driven by their own dedicated
	// TestAzure*EndToEnd functions below instead of runLabEndToEnd.
	"azure-001": {},
	"azure-002": {},
}

// TestAllLabsCovered catches a lab added to labs/ without a matching
// entry (and variant list) in endToEndLabs above — a plain os.ReadDir
// check, so it runs even without BREAKLAB_INTEGRATION=1.
func TestAllLabsCovered(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "..", "labs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, ok := endToEndLabs[e.Name()]; !ok {
			t.Errorf("lab %q has no entry in endToEndLabs — add one so it gets an end-to-end check", e.Name())
		}
	}
	for lab := range endToEndLabs {
		if _, err := os.Stat(filepath.Join("..", "..", "labs", lab)); err != nil {
			t.Errorf("endToEndLabs has %q but labs/%s doesn't exist — remove the stale entry", lab, lab)
		}
	}
}

// runLabEndToEnd needs a running Docker daemon (and, for the two k3s
// labs, a privileged container), so it only runs with
// BREAKLAB_INTEGRATION=1. For each of lab's variants it verifies: the
// variant really breaks the lab (check fails), that variant's own
// solution script really fixes it (check passes, and reports back the
// same variant that was forced), and Reset reapplies the same breakage
// (check fails again).
func runLabEndToEnd(t *testing.T, lab string, variants []string) {
	t.Helper()
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}

	for _, variant := range variants {
		t.Run(variant, func(t *testing.T) {
			t.Setenv("BREAKLAB_FORCE_VARIANT", variant)

			ctx := context.Background()
			d := docker.New()
			r := New(d)
			r.LabsDir = filepath.Join("..", "..", "labs")
			r.RoadmapPath = filepath.Join("..", "..", "roadmap.yaml")
			r.SessionPath = filepath.Join(t.TempDir(), "session.json")
			r.HistoryPath = filepath.Join(t.TempDir(), "history.json")
			r.ProgressPath = filepath.Join(t.TempDir(), "progress.json")
			r.Out = os.Stderr
			t.Cleanup(func() { _ = r.Stop(context.Background()) })

			if err := r.Start(ctx, lab, true); err != nil {
				t.Fatal(err)
			}

			res, err := r.Check(ctx)
			if err != nil || res.Passed {
				t.Fatalf("broken lab should fail its check: %+v, %v", res, err)
			}

			container := "breaklab-" + lab
			solution := filepath.Join(r.LabsDir, lab, "solutions", variant+".sh")
			if err := d.Copy(ctx, solution, container, "/tmp/solution.sh"); err != nil {
				t.Fatal(err)
			}
			// A handful of labs genuinely exhaust a resource (e.g. pids)
			// as part of their broken state, which can make docker exec
			// itself transiently fail to start — exactly what a real
			// operator would just retry, so the harness does too.
			var out []byte
			var execErr error
			for attempt := 0; attempt < 5; attempt++ {
				out, execErr = d.Exec(ctx, container, "bash", "/tmp/solution.sh")
				if execErr == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			if execErr != nil {
				t.Fatalf("solution failed: %v\n%s", execErr, out)
			}

			res, err = r.Check(ctx)
			if err != nil || !res.Passed || res.Variant != variant {
				t.Fatalf("fixed lab should pass: %+v, %v", res, err)
			}

			if err := r.Reset(ctx); err != nil {
				t.Fatal(err)
			}
			if res, _ := r.Check(ctx); res.Passed {
				t.Fatal("reset lab should be broken again")
			}
		})
	}
}

func TestAlertNeverFiresEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "alert-never-fires", endToEndLabs["alert-never-fires"])
}

func TestBackendDependencyDownEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "backend-dependency-down", endToEndLabs["backend-dependency-down"])
}

func TestBindMountShadowEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "bind-mount-shadow", endToEndLabs["bind-mount-shadow"])
}

func TestBisectTheBreakEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "bisect-the-break", endToEndLabs["bisect-the-break"])
}

func TestBlindAfterRotationEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "blind-after-rotation", endToEndLabs["blind-after-rotation"])
}

func TestBloatedDataEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "bloated-data", endToEndLabs["bloated-data"])
}

func TestCacheDependencyDownEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "cache-dependency-down", endToEndLabs["cache-dependency-down"])
}

func TestClockSkewTimelineEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "clock-skew-timeline", endToEndLabs["clock-skew-timeline"])
}

func TestCronBlindspotEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "cron-blindspot", endToEndLabs["cron-blindspot"])
}

func TestCrashLoopWorkerEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "crash-loop-worker", endToEndLabs["crash-loop-worker"])
}

// TestDeploymentNeverReadyEndToEnd needs a real single-node k3s control
// plane (privileged container), so it's much slower than most of this
// suite.
func TestDeletedFileStillOpenEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "deleted-file-still-open", endToEndLabs["deleted-file-still-open"])
}

func TestDeploymentNeverReadyEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "deployment-never-ready", endToEndLabs["deployment-never-ready"])
}

func TestDriftedConfigEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "drifted-config", endToEndLabs["drifted-config"])
}

func TestFakeGreenEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "fake-green", endToEndLabs["fake-green"])
}

func TestGitLostWorkEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "git-lost-work", endToEndLabs["git-lost-work"])
}

// TestHalfShippedEndToEnd drives the capstone lab directly rather than
// through runLabEndToEnd: Start refuses a capstone lab until every
// regular stage is complete (runner.go's checkCapstoneUnlocked), so this
// seeds a progress.json where that's already true before the usual
// break/fix/reset sequence.
func TestHalfShippedEndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}

	for _, variant := range endToEndLabs["half-shipped"] {
		t.Run(variant, func(t *testing.T) {
			t.Setenv("BREAKLAB_FORCE_VARIANT", variant)

			ctx := context.Background()
			d := docker.New()
			r := New(d)
			r.LabsDir = filepath.Join("..", "..", "labs")
			r.RoadmapPath = filepath.Join("..", "..", "roadmap.yaml")
			r.SessionPath = filepath.Join(t.TempDir(), "session.json")
			r.HistoryPath = filepath.Join(t.TempDir(), "history.json")
			r.ProgressPath = filepath.Join(t.TempDir(), "progress.json")
			r.Out = os.Stderr
			t.Cleanup(func() { _ = r.Stop(context.Background()) })

			seedAllRegularStagesSolved(t, r.ProgressPath, r.LabsDir, r.RoadmapPath)

			if err := r.Start(ctx, "half-shipped", true); err != nil {
				t.Fatal(err)
			}

			res, err := r.Check(ctx)
			if err != nil || res.Passed {
				t.Fatalf("broken lab should fail its check: %+v, %v", res, err)
			}

			container := "breaklab-half-shipped"
			solution := filepath.Join(r.LabsDir, "half-shipped", "solutions", variant+".sh")
			if err := d.Copy(ctx, solution, container, "/tmp/solution.sh"); err != nil {
				t.Fatal(err)
			}
			if out, err := d.Exec(ctx, container, "bash", "/tmp/solution.sh"); err != nil {
				t.Fatalf("solution failed: %v\n%s", err, out)
			}

			res, err = r.Check(ctx)
			if err != nil || !res.Passed || res.Variant != variant {
				t.Fatalf("fixed lab should pass: %+v, %v", res, err)
			}

			if err := r.Reset(ctx); err != nil {
				t.Fatal(err)
			}
			if res, _ := r.Check(ctx); res.Passed {
				t.Fatal("reset lab should be broken again")
			}
		})
	}
}

// TestK8sBlindSpotEndToEnd drives this capstone lab directly rather than
// through runLabEndToEnd, for the same reason as TestHalfShippedEndToEnd:
// Start refuses a capstone lab until every regular stage is complete.
func TestK8sBlindSpotEndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}

	for _, variant := range endToEndLabs["k8s-blind-spot"] {
		t.Run(variant, func(t *testing.T) {
			t.Setenv("BREAKLAB_FORCE_VARIANT", variant)

			ctx := context.Background()
			d := docker.New()
			r := New(d)
			r.LabsDir = filepath.Join("..", "..", "labs")
			r.RoadmapPath = filepath.Join("..", "..", "roadmap.yaml")
			r.SessionPath = filepath.Join(t.TempDir(), "session.json")
			r.HistoryPath = filepath.Join(t.TempDir(), "history.json")
			r.ProgressPath = filepath.Join(t.TempDir(), "progress.json")
			r.Out = os.Stderr
			t.Cleanup(func() { _ = r.Stop(context.Background()) })

			seedAllRegularStagesSolved(t, r.ProgressPath, r.LabsDir, r.RoadmapPath)

			if err := r.Start(ctx, "k8s-blind-spot", true); err != nil {
				t.Fatal(err)
			}

			res, err := r.Check(ctx)
			if err != nil || res.Passed {
				t.Fatalf("broken lab should fail its check: %+v, %v", res, err)
			}

			container := "breaklab-k8s-blind-spot"
			solution := filepath.Join(r.LabsDir, "k8s-blind-spot", "solutions", variant+".sh")
			if err := d.Copy(ctx, solution, container, "/tmp/solution.sh"); err != nil {
				t.Fatal(err)
			}
			var out []byte
			var execErr error
			for attempt := 0; attempt < 5; attempt++ {
				out, execErr = d.Exec(ctx, container, "bash", "/tmp/solution.sh")
				if execErr == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			if execErr != nil {
				t.Fatalf("solution failed: %v\n%s", execErr, out)
			}

			res, err = r.Check(ctx)
			if err != nil || !res.Passed || res.Variant != variant {
				t.Fatalf("fixed lab should pass: %+v, %v", res, err)
			}

			if err := r.Reset(ctx); err != nil {
				t.Fatal(err)
			}
			if res, _ := r.Check(ctx); res.Passed {
				t.Fatal("reset lab should be broken again")
			}
		})
	}
}

// TestPhantomReleaseEndToEnd drives this capstone lab directly rather
// than through runLabEndToEnd, for the same reason as
// TestHalfShippedEndToEnd: Start refuses a capstone lab until every
// regular stage is complete.
func TestPhantomReleaseEndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}

	for _, variant := range endToEndLabs["phantom-release"] {
		t.Run(variant, func(t *testing.T) {
			t.Setenv("BREAKLAB_FORCE_VARIANT", variant)

			ctx := context.Background()
			d := docker.New()
			r := New(d)
			r.LabsDir = filepath.Join("..", "..", "labs")
			r.RoadmapPath = filepath.Join("..", "..", "roadmap.yaml")
			r.SessionPath = filepath.Join(t.TempDir(), "session.json")
			r.HistoryPath = filepath.Join(t.TempDir(), "history.json")
			r.ProgressPath = filepath.Join(t.TempDir(), "progress.json")
			r.Out = os.Stderr
			t.Cleanup(func() { _ = r.Stop(context.Background()) })

			seedAllRegularStagesSolved(t, r.ProgressPath, r.LabsDir, r.RoadmapPath)

			if err := r.Start(ctx, "phantom-release", true); err != nil {
				t.Fatal(err)
			}

			res, err := r.Check(ctx)
			if err != nil || res.Passed {
				t.Fatalf("broken lab should fail its check: %+v, %v", res, err)
			}

			container := "breaklab-phantom-release"
			solution := filepath.Join(r.LabsDir, "phantom-release", "solutions", variant+".sh")
			if err := d.Copy(ctx, solution, container, "/tmp/solution.sh"); err != nil {
				t.Fatal(err)
			}
			if out, err := d.Exec(ctx, container, "bash", "/tmp/solution.sh"); err != nil {
				t.Fatalf("solution failed: %v\n%s", err, out)
			}

			res, err = r.Check(ctx)
			if err != nil || !res.Passed || res.Variant != variant {
				t.Fatalf("fixed lab should pass: %+v, %v", res, err)
			}

			if err := r.Reset(ctx); err != nil {
				t.Fatal(err)
			}
			if res, _ := r.Check(ctx); res.Passed {
				t.Fatal("reset lab should be broken again")
			}
		})
	}
}

// seedAllRegularStagesSolved writes one solved attempt for every lab in
// every regular (non-capstone) stage, unlocking the capstone for a test
// that needs to call Start on it.
func seedAllRegularStagesSolved(t *testing.T, progressPath, labsDir, roadmapPath string) {
	t.Helper()
	cat, err := catalog.LoadCatalog(labsDir, roadmapPath)
	if err != nil {
		t.Fatal(err)
	}
	pf := &progress.ProgressFile{Version: 1}
	for _, group := range cat.RegularStages() {
		for _, lab := range group.Labs {
			pf.Attempts = append(pf.Attempts, progress.Attempt{Lab: lab.ID, Status: progress.StatusSolved})
		}
	}
	if err := progress.Save(progressPath, pf); err != nil {
		t.Fatal(err)
	}
}

func TestImagePullBackoffEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "image-pull-backoff", endToEndLabs["image-pull-backoff"])
}

func TestImpossibleSchedulingRequestEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "impossible-scheduling-request", endToEndLabs["impossible-scheduling-request"])
}

func TestInfraEnvOverrideEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "infra-env-override", endToEndLabs["infra-env-override"])
}

func TestInfraStuckLockEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "infra-stuck-lock", endToEndLabs["infra-stuck-lock"])
}

// TestK8sa001EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sa001EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sa-001", endToEndLabs["k8sa-001"])
}

// TestK8sa002EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sa002EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sa-002", endToEndLabs["k8sa-002"])
}

// TestK8sa003EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sa003EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sa-003", endToEndLabs["k8sa-003"])
}

// TestK8sd001EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sd001EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-001", endToEndLabs["k8sd-001"])
}

// TestK8sd002EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sd002EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-002", endToEndLabs["k8sd-002"])
}

// TestK8sd003EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sd003EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-003", endToEndLabs["k8sd-003"])
}

// TestK8sd004EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sd004EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-004", endToEndLabs["k8sd-004"])
}

// TestK8sd005EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sd005EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-005", endToEndLabs["k8sd-005"])
}

// TestK8sd006EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sd006EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-006", endToEndLabs["k8sd-006"])
}

// TestK8sd007EndToEnd needs a real single-node k3s control plane
// (privileged container), plus a real ArgoCD core install and git
// daemon, so it's much slower than most of this suite.
func TestK8sd007EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-007", endToEndLabs["k8sd-007"])
}

// TestK8sd008EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sd008EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sd-008", endToEndLabs["k8sd-008"])
}

// TestK8sf001EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sf001EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sf-001", endToEndLabs["k8sf-001"])
}

// TestK8sf002EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sf002EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sf-002", endToEndLabs["k8sf-002"])
}

// TestK8sf003EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sf003EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sf-003", endToEndLabs["k8sf-003"])
}

// TestK8sf004EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sf004EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sf-004", endToEndLabs["k8sf-004"])
}

// TestK8sf005EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sf005EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sf-005", endToEndLabs["k8sf-005"])
}

// TestK8sf006EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8sf006EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8sf-006", endToEndLabs["k8sf-006"])
}

// TestK8ss001EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8ss001EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8ss-001", endToEndLabs["k8ss-001"])
}

// TestK8ss002EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8ss002EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8ss-002", endToEndLabs["k8ss-002"])
}

// TestK8ss003EndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestK8ss003EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8ss-003", endToEndLabs["k8ss-003"])
}

// TestK8ss004EndToEnd needs a real single-node k3s control plane
// (privileged container), plus a real Kyverno install, so it's much
// slower than most of this suite.
func TestK8ss004EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "k8ss-004", endToEndLabs["k8ss-004"])
}

func TestLeakedSecretEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "leaked-secret", endToEndLabs["leaked-secret"])
}

func TestMetricRenameGapEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "metric-rename-gap", endToEndLabs["metric-rename-gap"])
}

func TestNetworkPolicyBlockEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "network-policy-block", endToEndLabs["network-policy-block"])
}

func TestOOMKilledEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "oom-killed", endToEndLabs["oom-killed"])
}

func TestPathShadowEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "path-shadow", endToEndLabs["path-shadow"])
}

func TestPidsCeilingEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "pids-ceiling", endToEndLabs["pids-ceiling"])
}

func TestPinnedModuleRegressionEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "pinned-module-regression", endToEndLabs["pinned-module-regression"])
}

func TestProxyMazeEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "proxy-maze", endToEndLabs["proxy-maze"])
}

func TestPVCBadStorageClassEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "pvc-bad-storageclass", endToEndLabs["pvc-bad-storageclass"])
}

func TestQueueDrainEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "queue-drain", endToEndLabs["queue-drain"])
}

func TestQuietShipperEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "quiet-shipper", endToEndLabs["quiet-shipper"])
}

func TestReadyButNotReadyEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "ready-but-not-ready", endToEndLabs["ready-but-not-ready"])
}

func TestSelectorMismatchEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "selector-mismatch", endToEndLabs["selector-mismatch"])
}

func TestShipperBufferDropEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "shipper-buffer-drop", endToEndLabs["shipper-buffer-drop"])
}

func TestSilentPipelineEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "silent-pipeline", endToEndLabs["silent-pipeline"])
}

func TestSilentStdoutEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "silent-stdout", endToEndLabs["silent-stdout"])
}

func TestStaleCacheRestoreEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "stale-cache-restore", endToEndLabs["stale-cache-restore"])
}

// TestStaleConfigEndToEnd needs a real single-node k3s control plane
// (privileged container), so it's much slower than most of this suite.
func TestStaleConfigEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "stale-config", endToEndLabs["stale-config"])
}

func TestStaleConfigMapEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "stale-configmap", endToEndLabs["stale-configmap"])
}

func TestStaleStashEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "stale-stash", endToEndLabs["stale-stash"])
}

func TestStaleSymlinkEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "stale-symlink", endToEndLabs["stale-symlink"])
}

func TestStuckLockEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "stuck-lock", endToEndLabs["stuck-lock"])
}

func TestSwallowedSignalEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "swallowed-signal", endToEndLabs["swallowed-signal"])
}

func TestUmaskMismatchEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "umask-mismatch", endToEndLabs["umask-mismatch"])
}

func TestUnexpandedTemplateEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "unexpanded-template", endToEndLabs["unexpanded-template"])
}

func TestUnmanagedResourceEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "unmanaged-resource", endToEndLabs["unmanaged-resource"])
}

func TestWebDownEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "web-down", endToEndLabs["web-down"])
}

func TestWrongArtifactPromotedEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "wrong-artifact-promoted", endToEndLabs["wrong-artifact-promoted"])
}

func TestZombiePileupEndToEnd(t *testing.T) {
	runLabEndToEnd(t, "zombie-pileup", endToEndLabs["zombie-pileup"])
}

// newIntegrationRunner builds a *Runner against the real labs/ and
// roadmap.yaml with every state path under t.TempDir(), for the
// standalone-lab (lnx-*) end-to-end tests below, which don't fit
// runLabEndToEnd's variant/solutions-script shape.
func newIntegrationRunner(t *testing.T) *Runner {
	t.Helper()
	d := docker.New()
	r := New(d)
	r.LabsDir = filepath.Join("..", "..", "labs")
	r.RoadmapPath = filepath.Join("..", "..", "roadmap.yaml")
	r.SessionPath = filepath.Join(t.TempDir(), "session.json")
	r.HistoryPath = filepath.Join(t.TempDir(), "history.json")
	r.ProgressPath = filepath.Join(t.TempDir(), "progress.json")
	r.Out = os.Stderr
	t.Cleanup(func() { _ = r.Stop(context.Background()) })
	return r
}

// TestLnx001EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestLnx001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
DISTRIBUTION=$(. /etc/os-release; echo "$ID")
KERNEL=$(uname -r)
HOSTNAME=$(hostname)
USER=$(whoami)
SHELL=$SHELL
WORKING_DIRECTORY=$(pwd)
ARCHITECTURE=$(uname -m)
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx003EndToEnd proves the Guided lab's starting state fails (no
// transcript yet) and the documented solution passes.
func TestLnx003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no transcript yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /etc && pwd >> /workspace/transcript.txt
cd /var/log && pwd >> /workspace/transcript.txt
cd /opt && pwd >> /workspace/transcript.txt
cd ~ && pwd >> /workspace/transcript.txt
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx005EndToEnd proves the Guided lab's starting state fails (the
// app layout doesn't exist yet) and the documented solution passes.
func TestLnx005EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-005", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-005"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no app layout yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
mkdir -p /opt/acme-app/{bin,config,data,logs}
touch /opt/acme-app/config/app.conf
touch /opt/acme-app/logs/app.log
cp /opt/fixtures/acme-app /opt/acme-app/bin/acme-app
chmod +x /opt/acme-app/bin/acme-app
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs003EndToEnd proves the Guided lab's starting state fails (no
// alert rule yet) and the documented solution passes.
func TestObs003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no alert rule yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /etc/alerts/rules.yaml
metric: error_rate
op: ">"
threshold: 5
EOF
alert-evaluator
cat /var/log/alerts.log
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs004EndToEnd proves the Guided lab's starting state fails (the
// aggregator is still empty) and the documented solution passes.
func TestObs004EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-004", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-004"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("aggregator still empty: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
ship-logs --source /var/log/app.log --destination /var/log/aggregator.log

echo "2026-10-20T00:01:00Z app new line shipped later" >> /var/log/app.log
ship-logs --source /var/log/app.log --destination /var/log/aggregator.log
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs005EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestObs005EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-005", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-005"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
pid=$(cat /var/run/reader.pid)
fd_target=""
for fd in /proc/"$pid"/fd/*; do
  tgt=$(readlink -f "$fd")
  case "$tgt" in
    /var/log/app.log*) fd_target="$tgt" ;;
  esac
done
current_inode=$(stat -c %i /var/log/app.log)
fd_inode=$(stat -c %i "$fd_target")
inodes_match=no
[ "$current_inode" = "$fd_inode" ] && inodes_match=yes

cat <<EOF > /workspace/answers.env
READER_STALE=yes
READER_FD_PATH=$fd_target
INODES_MATCH=$inodes_match
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs006EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestObs006EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-006", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-006"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
panel_metric=$(sed -n 's/^metric: *//p' /etc/dashboards/overview.yaml)
exported_metric=disk_free_percent
render-dashboard
rendered=$(cat /var/lib/dashboards/overview.rendered)

cat <<EOF > /workspace/answers.env
PANEL_METRIC=$panel_metric
EXPORTED_METRIC=$exported_metric
RENDER_RESULT=$rendered
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs007EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestObs007EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-007", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-007"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
real_now=$(date +%s)
log-event-a "check"
log-event-b "check"

ts_a=$(tail -n 1 /var/log/host-a.log | awk '{print $2}')
ts_b=$(tail -n 1 /var/log/host-b.log | awk '{print $2}')

skew_a=$(( ts_a - real_now )); [ "$skew_a" -lt 0 ] && skew_a=$(( -skew_a ))
skew_b=$(( ts_b - real_now )); [ "$skew_b" -lt 0 ] && skew_b=$(( -skew_b ))

if [ "$skew_a" -gt "$skew_b" ]; then
  drifted=host-a; skew=$skew_a
else
  drifted=host-b; skew=$skew_b
fi

cat <<EOF > /workspace/answers.env
DRIFTED_HOST=$drifted
SKEW_SECONDS=$skew
IS_REAL_ORDERING_ISSUE=no
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx010EndToEnd proves the Practice lab's starting state fails (no
// pipeline written yet), the documented solution passes, and — per this
// lab's own constraint — the answer stays correct after the input log
// is reshuffled into a different line order.
func TestLnx010EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-010", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-010"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no pipeline written yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'SCRIPT' > /workspace/top-500-ip.sh
#!/usr/bin/env bash
awk '$5 == 500 {print $1}' /workspace/access.log | sort | uniq -c | sort -rn | head -1 | awk '{print $1, $2}'
SCRIPT
chmod +x /workspace/top-500-ip.sh
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass (including the reshuffle re-check the lab's own check.sh performs): %+v, %v", res, err)
	}
}

// TestLnx015EndToEnd proves lnx-015's multi-check contract end to end:
// the starting state fails overall (app-can-read fails even though
// permissions-safe already happens to pass — the config starts
// root-owned 600, which is restrictive but unreadable by appuser), the
// least-privilege fix passes both checks, and — the whole reason this
// lab exists as the multi-check proof — a chmod 777 shortcut would make
// app-can-read pass while permissions-safe correctly fails it.
func TestLnx015EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-015", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-015"

	res, err := r.Check(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed {
		t.Fatalf("broken lab should fail overall: %+v", res)
	}
	resultByID := map[string]bool{}
	for _, cr := range res.Results {
		resultByID[cr.ID] = cr.Passed
	}
	if resultByID["app-can-read"] {
		t.Errorf("app-can-read should fail before the fix: %+v", res.Results)
	}
	if !resultByID["permissions-safe"] {
		t.Errorf("permissions-safe should already pass before the fix (the broken state is restrictive, not unsafe): %+v", res.Results)
	}

	fix := `set -e
chgrp appuser /srv/app/config.yaml
chmod 640 /srv/app/config.yaml
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", fix); err != nil {
		t.Fatalf("applying fix: %v\n%s", err, out)
	}

	res, err = r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass both checks: %+v, %v", res, err)
	}
	for _, cr := range res.Results {
		if !cr.Passed {
			t.Errorf("check %q should pass after the fix: %+v", cr.ID, cr)
		}
	}

	// Now prove the shortcut this lab exists to block: chmod 777 makes
	// app-can-read pass, but permissions-safe must still fail it.
	shortcut := "chmod 777 /srv/app/config.yaml"
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", shortcut); err != nil {
		t.Fatalf("applying shortcut: %v\n%s", err, out)
	}
	res, err = r.Check(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed {
		t.Fatalf("chmod 777 shortcut should still fail overall: %+v", res)
	}
	resultByID = map[string]bool{}
	for _, cr := range res.Results {
		resultByID[cr.ID] = cr.Passed
	}
	if !resultByID["app-can-read"] {
		t.Errorf("app-can-read should pass after chmod 777 (that's why it's not sufficient on its own): %+v", res.Results)
	}
	if resultByID["permissions-safe"] {
		t.Errorf("permissions-safe should reject the chmod 777 shortcut: %+v", res.Results)
	}
}

// TestLnx002EndToEnd proves the Practice lab's starting state fails (no
// answers.env, no fixture files created yet) and the documented solution
// passes.
func TestLnx002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-002", true); err != nil {
		t.Fatal(err)
	}

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("nothing done yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
HIDDEN_FILES_OPTION=-a
NESTED_MKDIR_OPTION=-p
INTERACTIVE_COPY_OPTION=-i
EOF
touch /workspace/.hidden-check
mkdir -p /workspace/deep/nested/tree
cp -i /opt/fixtures/readme.txt /workspace/fixture-copy.txt
`
	if out, err := r.Docker.Exec(ctx, "breaklab-lnx-002", "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx004EndToEnd proves the Practice lab's starting state fails (no
// answer.txt yet) and the documented solution — a content search that
// finds the one buried log containing the transaction ID — passes.
func TestLnx004EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-004", true); err != nil {
		t.Fatal(err)
	}

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answer.txt yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
echo "/var/log/app/archive/2024-01/host-03/trace.log" > /workspace/answer.txt
`
	if out, err := r.Docker.Exec(ctx, "breaklab-lnx-004", "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx006EndToEnd proves the break-fix lab's starting state (three
// files seeded in the wrong directories under /opt/acme-app) fails the
// check, and moving each file to its correct directory passes.
func TestLnx006EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-006", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-006"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("broken layout: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
mv /opt/acme-app/config/acme-app /opt/acme-app/bin/acme-app
chmod +x /opt/acme-app/bin/acme-app
mv /opt/acme-app/data/app.conf /opt/acme-app/config/app.conf
mv /opt/acme-app/bin/data.db /opt/acme-app/data/data.db
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit001EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestGit001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
cat <<EOF > answers.env
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
HEAD_COMMIT_SHORT_SHA=$(git rev-parse --short HEAD)
OTHER_BRANCH_NAME=$(git branch --format="%(refname:short)" | grep -v -x "$(git rev-parse --abbrev-ref HEAD)")
MODIFIED_FILE_COUNT=$(git status --porcelain | grep -c "^ M ")
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit002EndToEnd proves the Guided lab's starting state fails (the
// timeout edit is still uncommitted) and the documented solution —
// staging and committing only config.yaml, leaving scratch-notes.txt
// untracked — passes.
func TestGit002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("nothing committed yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git add config.yaml
git commit -m "Increase API timeout to 30s"
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit003EndToEnd proves the Practice lab's starting state fails
// (the bugfix and the unrelated debug print are both still uncommitted)
// and the documented solution — committing each file separately —
// passes.
func TestGit003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("nothing committed yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git add app.py
git commit -m "Fix off-by-one in get_last_item"
git add utils.py
git commit -m "Add debug print to utils.py"
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit004EndToEnd proves the Guided lab's starting state fails (no
// feature-logging branch yet) and the documented solution — branching,
// committing logging.md, and switching back to main — passes.
func TestGit004EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-004", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-004"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no feature-logging branch yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git checkout -b feature-logging
printf '# Logging plan\n\nAdd structured logging to the service.\n' > logging.md
git add logging.md
git commit -m "Add logging plan"
git checkout main
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit005EndToEnd proves the Practice lab's starting state fails
// (main still sits at the mistaken tip) and the documented solution —
// branch off the mistaken commits, then reset main back to
// known-good-base — passes.
func TestGit005EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-005", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-005"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("main still at the mistaken tip: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git branch feature/csv-export main
git reset --hard known-good-base
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit006EndToEnd proves the Guided lab's starting state fails (no
// merge commit yet) and the documented solution — merging
// feature/signup-form into main, which can't fast-forward since both
// branches moved — passes.
func TestGit006EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-006", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-006"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("not merged yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git checkout main
git merge --no-edit feature/signup-form
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit007EndToEnd proves the break-fix lab's starting state fails (no
// merge commit yet) and the documented solution — merging
// feature/raise-timeout, resolving the real conflict in config.yaml to
// the agreed content, and completing the merge commit — passes.
func TestGit007EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-007", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-007"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("not merged yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git checkout main
git merge feature/raise-timeout || true
cat <<EOF > config.yaml
service: acme-api
port: 8080
timeout: 60
EOF
git add config.yaml
git commit --no-edit
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit008EndToEnd proves the Guided lab's starting state fails
// (local main hasn't pushed anything back to origin yet) and the
// documented solution — fetch, pull the teammate's commit in, commit
// locally, and push — passes.
func TestGit008EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-008", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-008"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("not pushed yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git fetch origin
git pull origin main
cat <<EOF >> CHANGELOG.md
## Changelog
- Added greeting module
- Pulled in logging helper from teammate
EOF
git add CHANGELOG.md
git commit -m "Add changelog"
git push origin main
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx007EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — the 3 keys that
// actually differ between config.old and config.new — passes.
func TestLnx007EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-007", true); err != nil {
		t.Fatal(err)
	}

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
DATABASE_HOST_OLD=db-old.internal
DATABASE_HOST_NEW=db-new.internal
LOG_LEVEL_OLD=debug
LOG_LEVEL_NEW=info
TIMEOUT_OLD=30
TIMEOUT_NEW=60
EOF
`
	if out, err := r.Docker.Exec(ctx, "breaklab-lnx-007", "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx008EndToEnd proves lnx-008's safety guardrail end to end: the
// starting state fails (the live config still points at the old
// backend), a blind find-and-replace across the whole /opt/app tree
// still fails the check (because it corrupts the archived/historical
// files that are supposed to keep mentioning the old hostname), and the
// documented solution — editing only the live config — passes while
// leaving the archived file's content untouched.
func TestLnx008EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-008", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-008"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("stale backend: check should fail, got %+v, %v", res, err)
	}

	// Prove the shortcut this lab exists to block: a blind recursive
	// find-and-replace "fixes" the live config but also corrupts the
	// archived release snapshot and historical docs — the check must
	// still fail.
	blindReplace := `set -e
grep -rlF "database-old.internal" /opt/app | xargs sed -i 's/database-old\.internal/database.internal/'
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", blindReplace); err != nil {
		t.Fatalf("applying blind replace: %v\n%s", err, out)
	}
	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("blind find-and-replace across the whole tree should still fail the check: %+v, %v", res, err)
	}
	archived, err := r.Docker.Exec(ctx, container, "cat", "/opt/app/releases/v1/config/app.conf")
	if err != nil {
		t.Fatalf("reading archived file: %v", err)
	}
	if strings.Contains(string(archived), "database-old.internal") {
		t.Fatalf("test setup error: blind replace didn't actually corrupt the archived file as expected:\n%s", archived)
	}

	// Reset to the broken starting state and apply the real fix: edit
	// only the live config.
	if err := r.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if err := r.Start(ctx, "lnx-008", true); err != nil {
		t.Fatal(err)
	}

	fix := `set -e
sed -i 's/database-old\.internal/database.internal/' /opt/app/config/app.conf
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", fix); err != nil {
		t.Fatalf("applying fix: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}

	archivedAfter, err := r.Docker.Exec(ctx, container, "cat", "/opt/app/releases/v1/config/app.conf")
	if err != nil {
		t.Fatalf("reading archived file after fix: %v", err)
	}
	if !strings.Contains(string(archivedAfter), "database-old.internal") {
		t.Fatalf("archived file should still contain database-old.internal after the real fix, got:\n%s", archivedAfter)
	}
}

// TestLnx009EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestLnx009EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-009", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-009"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
TOTAL_REQUESTS=$(wc -l < /workspace/access.log)
FIRST_TIMESTAMP=$(head -1 /workspace/access.log | cut -d' ' -f2)
LAST_TIMESTAMP=$(tail -1 /workspace/access.log | cut -d' ' -f2)
COUNT_2XX=$(cut -d' ' -f5 /workspace/access.log | grep -c '^2')
COUNT_4XX=$(cut -d' ' -f5 /workspace/access.log | grep -c '^4')
COUNT_5XX=$(cut -d' ' -f5 /workspace/access.log | grep -c '^5')
MOST_FREQUENT_PATH=$(cut -d' ' -f4 /workspace/access.log | sort | uniq -c | sort -rn | head -1 | awk '{print $2}')
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx011EndToEnd proves the Practice lab's starting state fails (no
// normalized.txt yet) and the documented solution passes.
func TestLnx011EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-011", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-011"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no normalized.txt yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
tr ',' '\n' < /workspace/hosts.txt \
  | tr 'A-Z' 'a-z' \
  | sed 's/^[ \t]*//; s/[ \t]*$//' \
  | sed '/^$/d' \
  | sort -u \
  > /workspace/normalized.txt
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx012EndToEnd proves the Guided lab's starting state fails (no
// developer/operator accounts yet) and the documented solution passes.
func TestLnx012EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-012", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-012"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no accounts yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
useradd -m -s /bin/bash developer
useradd -r -s /usr/sbin/nologin -g operator operator
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx013EndToEnd proves the Practice lab's starting state fails (no
// ops-team group yet, /srv/team not shared) and the documented solution
// passes.
func TestLnx013EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-013", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-013"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no ops-team sharing yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
groupadd ops-team
usermod -aG ops-team developer
usermod -aG ops-team operator
chgrp ops-team /srv/team
chmod 2770 /srv/team
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx014EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestLnx014EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-014", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-014"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
OWNER_ONLY_READABLE_BY_TESTER=no
GROUP_READABLE_READABLE_BY_TESTER=no
WORLD_READABLE_READABLE_BY_TESTER=yes
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx016EndToEnd proves the Practice lab's starting state fails (no
// /srv/team/newfile.txt yet) and the documented umask fix passes.
func TestLnx016EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-016", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-016"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no newfile.txt yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
umask 007
touch /srv/team/newfile.txt
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx017EndToEnd proves the break-fix lab's starting state fails
// (the runaway busy loop is still running, no answers.env yet) and the
// documented fix (identify it, record it, SIGTERM it) passes while the
// two healthy background workers stay untouched.
func TestLnx017EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-017", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-017"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("runaway still running, no answers.env: check should fail, got %+v, %v", res, err)
	}

	// The fix script's own text mentions "data-compactor"/"while true",
	// the very patterns check.sh's pgrep looks for — if that text sat
	// inline in a "bash -c <script>" invocation, that invoking process's
	// own argv would itself match pgrep -f and throw off the PID lookup.
	// Writing it to a file first and running the file keeps the pattern
	// text out of any live process's argv once pgrep actually runs.
	write := `cat <<'EOF' > /tmp/fix.sh
set -e
pid=$(pgrep -f data-compactor)
cat <<CFG > /workspace/answers.env
OFFENDING_PID=$pid
OFFENDING_PPID=1
OFFENDING_USER=root
OFFENDING_CMD="data-compactor bash -c while true; do :; done"
CFG
kill "$pid"
sleep 1
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", write); err != nil {
		t.Fatalf("writing fix script: %v\n%s", err, out)
	}
	if out, err := r.Docker.Exec(ctx, container, "bash", "/tmp/fix.sh"); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx018EndToEnd proves the guided lab's starting state fails (no
// graceful marker yet, no answers.env) and the documented SIGTERM/SIGKILL
// sequence passes: the gracefully-terminated worker1 leaves its marker,
// the force-killed second instance leaves none.
func TestLnx018EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-018", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-018"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no graceful marker yet: check should fail, got %+v, %v", res, err)
	}

	// As in TestLnx017EndToEnd: keep "worker1" out of the invoking
	// process's own argv, or pgrep -f worker1 matches that invoking
	// process too and "kill $pid1" ends up signaling the test's own
	// shell along with the real target.
	write := `cat <<'EOF' > /tmp/fix.sh
set -e
pid1=$(pgrep -f worker1)
kill "$pid1"
sleep 1
/usr/local/bin/worker.sh /workspace/worker-forced.marker &
pid2=$!
sleep 0.2
kill -9 "$pid2"
sleep 1
cat <<CFG > /workspace/answers.env
GRACEFUL_MARKER_EXISTS=yes
FORCED_MARKER_EXISTS=no
CFG
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", write); err != nil {
		t.Fatalf("writing fix script: %v\n%s", err, out)
	}
	if out, err := r.Docker.Exec(ctx, container, "bash", "/tmp/fix.sh"); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx019EndToEnd proves the break-fix lab's starting state fails
// (start.sh always exits 0 even with APP_CONFIG_PATH unset) and the
// documented rewrite passes both the "must fail" and "must succeed"
// sub-checks.
func TestLnx019EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-019", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-019"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("broken start.sh: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /opt/app/start.sh
#!/usr/bin/env bash
set -euo pipefail

if [ -z "${APP_CONFIG_PATH:-}" ]; then
  echo "ERROR: APP_CONFIG_PATH is not set; cannot start" >&2
  exit 1
fi

if [ ! -f "$APP_CONFIG_PATH" ]; then
  echo "ERROR: configuration file not found" >&2
  exit 1
fi

if ! grep -q "service:" "$APP_CONFIG_PATH"; then
  echo "ERROR: configuration file is invalid" >&2
  exit 1
fi

echo "Service started successfully"
exit 0
EOF
chmod +x /opt/app/start.sh
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx020EndToEnd proves the guided lab's starting state fails (none
// of curl/ip/lsof/jq are installed on the minimal base image) and the
// documented apt-get install passes.
func TestLnx020EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-020", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-020"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("minimal image: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
apt-get update
apt-get install -y curl iproute2 lsof jq
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx021EndToEnd proves the practice lab's starting state fails (no
// answers.env yet) and the documented du/find-derived solution passes.
func TestLnx021EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-021", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-021"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
LARGEST_DIR=/data/cache
LARGEST_FILE=/data/cache/big-dump.bin
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx022EndToEnd proves the incident lab's starting state fails (the
// debug directory is still full of junk) and that clearing only the
// junk, plus proving a new log can be written, passes.
func TestLnx022EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-022", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-022"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("junk data still present: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
rm -rf /var/log/app/debug
touch /var/log/app/new.log
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx023EndToEnd proves the practice lab's starting state fails (no
// answers.env yet) and the documented cross-log correlation passes.
func TestLnx023EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-023", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-023"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
FAILING_DEPENDENCY=payments
REQUEST_ID=req-8841
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx024EndToEnd proves the project lab's starting state fails (no
// health.sh yet), then — per this lab's own validation spec — explicitly
// proves health.sh itself (not just the lab's check.sh) returns non-zero
// against an unhealthy fixture, before proving the lab's own check.sh
// (which runs its own healthy/break/unhealthy dance) passes.
func TestLnx024EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-024", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-024"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no health.sh yet: check should fail, got %+v, %v", res, err)
	}

	healthScript := `#!/usr/bin/env bash
set -uo pipefail

PROCESS_PATTERN="worker.sh"
PORT=80
HEALTH_URL="http://localhost:${PORT}/healthz"
DISK_PATH="/var/log/app"
DISK_THRESHOLD_MB=100

if ! pgrep -f "$PROCESS_PATTERN" >/dev/null 2>&1; then
  echo "UNHEALTHY: required process ($PROCESS_PATTERN) is not running"
  exit 1
fi

if ! (exec 3<>"/dev/tcp/127.0.0.1/${PORT}") 2>/dev/null; then
  echo "UNHEALTHY: nothing is listening on port ${PORT}"
  exit 1
fi
exec 3>&- 2>/dev/null || true

code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 3 "$HEALTH_URL")
if [ "$code" != "200" ]; then
  echo "UNHEALTHY: $HEALTH_URL returned $code"
  exit 1
fi

disk_mb=$(du -sm "$DISK_PATH" 2>/dev/null | cut -f1)
disk_mb=${disk_mb:-0}
if [ "$disk_mb" -ge "$DISK_THRESHOLD_MB" ]; then
  echo "UNHEALTHY: $DISK_PATH is using ${disk_mb}MB (threshold ${DISK_THRESHOLD_MB}MB)"
  exit 1
fi

echo "HEALTHY: process up, port ${PORT} listening, ${HEALTH_URL} returned 200, disk ${disk_mb}MB/${DISK_THRESHOLD_MB}MB"
exit 0
`
	install := "set -e\ncat <<'SCRIPT' > /workspace/health.sh\n" + healthScript + "SCRIPT\nchmod +x /workspace/health.sh\n"
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", install); err != nil {
		t.Fatalf("installing health.sh: %v\n%s", err, out)
	}

	// Explicitly prove the unhealthy-fixture path returns non-zero from
	// health.sh itself, independent of the lab's own check.sh.
	if out, err := r.Docker.Exec(ctx, container, "bash", "/workspace/health.sh"); err != nil {
		t.Fatalf("health.sh should pass against the healthy fixture: %v\n%s", err, out)
	}
	if _, err := r.Docker.Exec(ctx, container, "pkill", "-f", "worker.sh"); err != nil {
		t.Fatalf("killing worker: %v", err)
	}
	time.Sleep(1 * time.Second)
	if out, err := r.Docker.Exec(ctx, container, "bash", "/workspace/health.sh"); err == nil {
		t.Fatalf("health.sh should fail once the worker process is gone, got: %s", out)
	}
	// Restore the worker so the lab's own check.sh sees a healthy
	// starting point when it runs its own healthy/unhealthy dance.
	if _, err := r.Docker.Exec(ctx, container, "bash", "-c", "nohup /opt/app/worker.sh >/dev/null 2>&1 & disown"); err != nil {
		t.Fatalf("restoring worker: %v", err)
	}
	time.Sleep(1 * time.Second)

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx025EndToEnd proves the break-fix lab's starting state fails (the
// seeded backup.sh has four defects) and that the documented fix both
// backs up correctly — including the space-named file — and, per this
// lab's own validation spec, still fails safely (non-zero, no false
// success message) when the source directory goes missing.
func TestLnx025EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "lnx-025", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-lnx-025"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("broken backup.sh: check should fail, got %+v, %v", res, err)
	}

	fix := `set -e
cat <<'SCRIPT' > /opt/app/backup.sh
#!/usr/bin/env bash
set -u
SRC=/data/to-backup
DEST=/opt/backups
STAMP=$(date +%Y%m%d%H%M%S)
STAGE="/tmp/backup-stage-$STAMP"
ARCHIVE="$DEST/backup-$STAMP.tar.gz"

mkdir -p "$STAGE" "$DEST"

for f in "$SRC"/*; do
  cp "$f" "$STAGE"/ || { echo "backup failed: could not copy $f" >&2; exit 1; }
done

if ! tar czf "$ARCHIVE" -C "$STAGE" .; then
  echo "backup failed: tar did not complete successfully" >&2
  exit 1
fi

echo "backup succeeded: $ARCHIVE"
SCRIPT
chmod +x /opt/app/backup.sh
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", fix); err != nil {
		t.Fatalf("applying fix: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestLnx900EndToEnd drives the Linux stage's final incident — the older
// multi-variant chaos-injection pattern — through all five of its seeded
// fault scenarios, exactly like TestWebDownEndToEnd above.
func TestLnx900EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "lnx-900", endToEndLabs["lnx-900"])
}

// TestNet001EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — every value read live
// from the running container rather than hard-coded — passes.
func TestNet001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "net-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-net-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
INTERFACE_NAME=$(ip route show default | awk '{print $5; exit}')
CONTAINER_IP=$(ip -4 -o addr show "$(ip route show default | awk '{print $5; exit}')" | awk '{print $4}' | cut -d/ -f1)
LISTENING_PORTS=$(ss -ltn | awk 'NR>1 {print $4}' | sed 's/.*://' | sort -un | paste -sd, -)
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestNet002EndToEnd proves the break-fix lab's starting state fails
// (a decoy process holds port 80 before nginx can bind it) and the
// documented fix — stop only the decoy, then start nginx — passes.
func TestNet002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "net-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-net-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("decoy holding port 80: check should fail, got %+v, %v", res, err)
	}

	// Writing the fix to a file first, then running the file, keeps the
	// "http.server 80" pattern text out of the invoking process's own
	// argv — otherwise "bash -c '...pkill -f \"http.server 80\"...'"
	// would itself match pkill -f and get signaled along with the real
	// decoy (same pitfall as TestLnx017EndToEnd/TestLnx018EndToEnd).
	write := `cat <<'EOF' > /tmp/fix.sh
set -e
pkill -f "http.server 80"
sleep 1
nginx -t
nginx
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", write); err != nil {
		t.Fatalf("writing fix script: %v\n%s", err, out)
	}
	if out, err := r.Docker.Exec(ctx, container, "bash", "/tmp/fix.sh"); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestNet003EndToEnd proves the Practice lab's starting state fails (no
// answers.env yet) and the documented classification of all four
// endpoints — healthy, dns-failure, refused, timeout, in that order —
// passes.
func TestNet003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "net-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-net-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
ENDPOINT_1=healthy
ENDPOINT_2=dns-failure
ENDPOINT_3=refused
ENDPOINT_4=timeout
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestNet004EndToEnd proves the Guided lab's starting state fails (no
// capture.pcap yet) and the documented tcpdump-while-curling sequence —
// a real capture showing both the request and the response on port
// 8005 — passes.
func TestNet004EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "net-004", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-net-004"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no capture.pcap yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
tcpdump -i lo -w /workspace/capture.pcap port 8005 &
TCPDUMP_PID=$!
sleep 1
curl http://localhost:8005/
sleep 1
kill "$TCPDUMP_PID"
sleep 1
cat <<EOF > /workspace/answers.env
DEST_PORT=8005
CAPTURE_HAS_REQUEST=yes
CAPTURE_HAS_RESPONSE=yes
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestNet005EndToEnd proves the break-fix lab's starting state fails (no
// /etc/hosts entry for api.internal, so the name doesn't resolve) and
// the documented fix — adding the missing /etc/hosts entry, never a
// hard-coded IP in the request — passes.
func TestNet005EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "net-005", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-net-005"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("api.internal unresolved: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
echo "127.0.0.1 api.internal" >> /etc/hosts
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestNet006EndToEnd proves the Practice lab's starting state fails (no
// request has reached the fixture server yet) and the documented curl
// call — method, both headers, and the JSON body all correct at once —
// passes.
func TestNet006EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "net-006", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-net-006"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no request made yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
curl -sS -X POST \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer test-token-fixture' \
  -d '{"action": "sync"}' \
  http://127.0.0.1:8080/
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestNet007EndToEnd proves the Practice lab's starting state fails (no
// answers.env yet) and the documented classification of all five
// endpoints — dns-failure, refused, timeout, 404, 500, in that order —
// passes.
func TestNet007EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "net-007", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-net-007"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
ENDPOINT_1=dns-failure
ENDPOINT_2=refused
ENDPOINT_3=timeout
ENDPOINT_4=404
ENDPOINT_5=500
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestWeb001EndToEnd proves the Guided lab's starting state fails (nginx
// isn't running and has no location block for /status/) and the
// documented fix — add the location block, validate with nginx -t, then
// start nginx — passes.
func TestWeb001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "web-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-web-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("nginx not configured/running yet: check should fail, got %+v, %v", res, err)
	}

	write := `cat <<'EOF' > /tmp/fix.sh
set -e
sed -i '/^    location \/ {$/i\
    location /status/ {\
        alias /var/www/status-site/;\
        index index.html;\
    }\
' /etc/nginx/sites-available/default
nginx -t
nginx
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", write); err != nil {
		t.Fatalf("writing fix script: %v\n%s", err, out)
	}
	if out, err := r.Docker.Exec(ctx, container, "bash", "/tmp/fix.sh"); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestWeb002EndToEnd proves the Guided lab's starting state fails (the
// backend is up on its own port but nginx has no proxy_pass to it yet,
// so nginx just serves its own default page) and the documented
// reverse-proxy config passes.
func TestWeb002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "web-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-web-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("nginx not yet proxying: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOT' > /etc/nginx/sites-available/default
server {
    listen 80 default_server;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:8090;
        proxy_set_header Host $host;
    }
}
EOT
nginx -t
nginx -s reload
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestWeb003EndToEnd proves the break-fix lab's starting state fails
// (nginx's proxy_pass points at a port nothing listens on, so every
// request 502s) and the documented fix — correcting proxy_pass to the
// backend's real port, then reloading nginx — passes.
func TestWeb003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "web-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-web-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("bad upstream port: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
sed -i 's/127.0.0.1:8098/127.0.0.1:8095/' /etc/nginx/sites-available/default
nginx -t
nginx -s reload
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestNet900EndToEnd drives the networking stage's final incident — the
// older multi-variant chaos-injection pattern — through all five of its
// seeded fault scenarios across the client -> nginx -> api -> database
// chain, exactly like TestWebDownEndToEnd/TestLnx900EndToEnd above.
func TestNet900EndToEnd(t *testing.T) {
	runLabEndToEnd(t, "net-900", endToEndLabs["net-900"])
}

// TestGit009EndToEnd proves the break-fix lab's starting state (origin's
// main already carries a teammate's commit the local clone never
// fetched, while the local clone has its own unpushed commit) makes a
// plain git push fail, and reconciling via rebase then pushing brings
// both commits together on origin.
func TestGit009EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-009", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-009"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("diverged histories: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git pull --rebase origin main
git push origin main
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit010EndToEnd proves the guided lab's three independent
// situations (a staged-but-uncommitted change, a bad commit that's only
// ever been local, and the same mistake already on a shared remote) all
// fail the check up front, and resolving each with its own command
// (restore, reset, revert respectively) passes.
func TestGit010EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-010", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-010"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("nothing fixed yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace/repo-staged
git restore --staged config.yaml
git restore config.yaml

cd /workspace/repo-private
git reset --hard HEAD~1

cd /workspace/repo-shared
git revert --no-edit HEAD

cat <<EOF > /workspace/answers.env
STAGED_CHANGE_COMMAND=restore
PRIVATE_COMMIT_COMMAND=reset
SHARED_COMMIT_COMMAND=revert
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit011EndToEnd proves the guided lab's starting state fails (no
// answers.env recording the commands used yet) and the documented
// stash-push-then-pop round trip — including the untracked file —
// passes, leaving the stash stack empty.
func TestGit011EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-011", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-011"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
git stash push -u -m "wip"
git stash pop

cat <<EOF > /workspace/answers.env
STASH_COMMAND=push
STASH_INCLUDE_UNTRACKED_FLAG=-u
POP_COMMAND=pop
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestGit012EndToEnd proves the practice lab's starting state fails (no
// v1.0.1 tag yet) and the documented solution — cherry-picking only the
// null-pointer hotfix onto release and tagging the result — passes.
func TestGit012EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "git-012", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-git-012"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no v1.0.1 tag yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /workspace
hotfix_sha=$(git log --format="%H %s" main | grep "Fix null-pointer crash in handler.py" | cut -d' ' -f1)
git checkout release
git cherry-pick "$hotfix_sha"
git tag -a v1.0.1 -m "Release 1.0.1: null-pointer crash fix"
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker001EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestDocker001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
process_count=$(ps aux --no-headers | wc -l)
live_hostname=$(hostname)

cat <<EOF > /workspace/answers.env
PROCESS_COUNT_SEEN=$process_count
KERNEL_SHARED=yes
HOSTNAME_SOURCE=container-runtime
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker002EndToEnd proves the Practice (reading-comprehension) lab's
// starting state fails (no answers.env yet) and the documented solution
// passes.
func TestDocker002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /workspace/answers.env
CACHE_MISS_ON_SOURCE_CHANGE=yes
CACHE_BUSTING_INSTRUCTION="COPY . ."
BUILDER_TOOLS_IN_FINAL_IMAGE=no
FINAL_STAGE_COPIES_FROM_BUILDER=dist,node_modules
SAME_RUN_DELETE_AVOIDS_BLOAT=yes
LATER_RUN_DELETE_AVOIDS_BLOAT=no
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker003EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestDocker003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
pid1_command=$(tr '\0' ' ' < /proc/1/cmdline | sed -e 's/[[:space:]]*$//')

cat <<EOF > /workspace/answers.env
PID1_COMMAND="$pid1_command"
IS_SHELL_FORM=yes
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker004EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes.
func TestDocker004EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-004", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-004"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
APP_ENV=$(printenv APP_ENV)
LOG_LEVEL=$(printenv LOG_LEVEL)
APP_VERSION=$(printenv APP_VERSION)
OVERRIDE_CANDIDATE=LOG_LEVEL
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker005EndToEnd proves the Guided lab's starting state fails (the
// Dockerfile's /data/welcome.txt is shadowed by the mounted volume, so
// /data starts empty) and the documented solution passes.
func TestDocker005EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-005", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-005"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("/data starts shadowed empty by the mounted volume: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat > /data/welcome.txt <<'WELCOME'
Welcome to the shared data volume.
WELCOME

data_dev=$(stat -c '%d' /data)
root_dev=$(stat -c '%d' /)

cat <<EOF2 > /workspace/answers.env
WELCOME_FILE_RESTORED=yes
MOUNT_IS_DISTINCT=$([ "$data_dev" != "$root_dev" ] && echo yes || echo no)
EOF2
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker006EndToEnd proves the break-fix lab's starting state (the
// status-page server bound to 127.0.0.1 only) fails its host-run check —
// the mapped host port can't reach a server listening on loopback only —
// and rebinding it to 0.0.0.0 passes, exactly like TestWebDownEndToEnd
// above.
func TestDocker006EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-006", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-006"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("server bound to 127.0.0.1 only: host check should fail, got %+v, %v", res, err)
	}

	// pkill -x matches only the process's exact short name ("python3"),
	// not its full argv — unlike pkill -f, which would also match this
	// very bash -c invocation, since its own script argument contains the
	// literal text "http.server 80" (the pattern a naive pkill -f would
	// search for) and so would kill itself before ever reaching the real
	// server process.
	solution := `set -e
sed -i 's/--bind 127.0.0.1/--bind 0.0.0.0/' /usr/local/bin/serve.sh
pkill -x python3 || true
sleep 1
nohup /usr/local/bin/serve.sh >/var/log/serve.log 2>&1 &
disown
sleep 1
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker007EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution passes. The check itself
// runs on the host (run_in: host) since docker logs can only be read from
// outside the container.
func TestDocker007EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-007", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-007"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
STDOUT_CAPTURED=yes
FILE_LOG_CAPTURED=no
LOG_FILE_PATH=/var/log/app/app.log
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker008EndToEnd proves the Guided lab's starting state fails (nginx's
// proxy_pass points at the wrong backend port, so the mapped host port
// returns a 502) and the documented solution — repointing proxy_pass at the
// backend's real port and reloading nginx — passes.
func TestDocker008EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-008", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-008"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("proxy_pass points at the wrong port: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
sed -i 's/5002/5001/' /etc/nginx/sites-available/default
nginx -t
nginx -s reload
sleep 1

cat <<EOF > /workspace/answers.env
BACKEND_REAL_PORT=5001
ORIGINAL_PROXY_PASS_PORT=5002
PROXY_FIXED=yes
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestDocker009EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — detecting the real cgroup
// version and reading the real pids/memory/CPU ceilings from it — passes.
func TestDocker009EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "docker-009", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-docker-009"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
  version=v2
  pids_max=$(cat /sys/fs/cgroup/pids.max)
  mem_max=$(cat /sys/fs/cgroup/memory.max)
  cpu_max=$(cat /sys/fs/cgroup/cpu.max)
else
  version=v1
  pids_max=$(cat /sys/fs/cgroup/pids/pids.max)
  mem_max=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)
  cpu_max=$(cat /sys/fs/cgroup/cpu/cpu.cfs_quota_us)
fi

cat <<EOF > /workspace/answers.env
CGROUP_VERSION=$version
PIDS_MAX=$pids_max
MEMORY_MAX=$mem_max
CPU_MAX="$cpu_max"
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd001EndToEnd proves the Practice (reading-comprehension) lab's
// starting state fails (no answers.env yet) and the documented solution
// passes.
func TestCicd001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /workspace/answers.env
TRIGGER_EVENTS=push,workflow_dispatch
PARALLEL_JOBS=lint,build
WAITING_JOB=test
WAITING_ON=build
CACHE_KEY_FILE=requirements.lock
MASKED_FAILURE_JOB=test
SWALLOWED_STEP="Run linter"
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd002EndToEnd proves the Guided lab's starting state fails
// (run.sh still chains stages with ';', so a failing test.sh doesn't stop
// deploy.sh) and the documented solution — chaining the stages with '&&'
// instead — passes.
func TestCicd002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("run.sh still uses ';': check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /opt/pipeline/run.sh
#!/usr/bin/env bash
echo "=== build ===" && /opt/pipeline/build.sh && \
echo "=== test ===" && /opt/pipeline/test.sh && \
echo "=== deploy ===" && /opt/pipeline/deploy.sh
EOF
chmod +x /opt/pipeline/run.sh
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd003EndToEnd proves the Guided lab's starting state fails (build
// and deploy disagree on the artifact's filename) and the documented
// solution — making deploy.sh discover whatever build.sh actually
// produced — passes.
func TestCicd003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no build/deploy run yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /opt/pipeline/deploy.sh
#!/usr/bin/env bash
set -euo pipefail

ARTIFACTS_DIR=/opt/pipeline/artifacts
APP_DIR=/opt/app

ARTIFACT="$(ls "$ARTIFACTS_DIR"/app-bundle-*.tar.gz 2>/dev/null | head -n1)"

if [ -z "$ARTIFACT" ] || [ ! -f "$ARTIFACT" ]; then
  echo "ERROR: no build artifact found in $ARTIFACTS_DIR" >&2
  exit 1
fi

rm -rf "$APP_DIR"
mkdir -p "$APP_DIR"
tar -xzf "$ARTIFACT" -C "$APP_DIR"

version="$(cat "$APP_DIR/VERSION")"
echo "$version" > /workspace/deployed-version.txt
echo "deployed $version"
EOF
chmod +x /opt/pipeline/deploy.sh

/opt/pipeline/build.sh
/opt/pipeline/deploy.sh
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd004EndToEnd proves the break-fix lab's starting state (a
// hardcoded secret literal compared against itself) fails the check, and
// the documented solution — reading DEPLOY_TOKEN from the environment and
// comparing its hash instead — passes.
func TestCicd004EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-004", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-004"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("hardcoded secret still in deploy.sh: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /opt/pipeline/deploy.sh
#!/usr/bin/env bash
set -euo pipefail

EXPECTED_HASH_FILE=/opt/pipeline/.deploy_token.sha256

if [ -z "${DEPLOY_TOKEN:-}" ]; then
  echo "DEPLOY DENIED: DEPLOY_TOKEN not set" >&2
  exit 1
fi

actual_hash="$(printf '%s' "$DEPLOY_TOKEN" | sha256sum | awk '{print $1}')"
expected_hash="$(cat "$EXPECTED_HASH_FILE")"

if [ "$actual_hash" = "$expected_hash" ]; then
  echo "DEPLOYED" > /workspace/deploy-status.txt
  echo "deploy authorized, shipping release"
else
  echo "DEPLOY DENIED: invalid token" >&2
  exit 1
fi
EOF
chmod +x /opt/pipeline/deploy.sh

export DEPLOY_TOKEN=brk-ci-7f2a9d41c8b3
/opt/pipeline/deploy.sh
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd005EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — recording the lock's
// dead PID and stale verdict — passes.
func TestCicd005EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-005", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-005"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
LOCK_PID=51234
PID_RUNNING=no
LOCK_STATUS=stale
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd006EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet, and cache-key still prints a static string regardless
// of requirements.lock's content) and the documented solution — deriving
// the key from a hash of requirements.lock instead — passes.
func TestCicd006EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-006", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-006"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /usr/local/bin/cache-key
#!/usr/bin/env bash
echo "deps-$(sha256sum /srv/app/requirements.lock | cut -c1-12)"
EOF
chmod +x /usr/local/bin/cache-key

cat <<EOF > /workspace/answers.env
ORIGINAL_KEY_SOURCE=static
CACHE_KEY_NOW_TRACKS_LOCKFILE=yes
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd007EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — identifying the
// symlink's real target instead of the newest-looking filename —
// passes.
func TestCicd007EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-007", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-007"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
PROMOTED_ARTIFACT=app-2.1.0.tar
WOULD_GUESSING_BY_NAME_BE_RIGHT=no
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac003EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — running infra-apply
// once with nothing changed, editing resources.yaml, then running it
// again — passes.
func TestIac003EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-003", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-003"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
infra-apply
sed -i 's/^value:.*/value: release-2026.10.2/' /infra/resources.yaml
infra-apply
cat <<EOF > /workspace/answers.env
FIRST_RUN_RESULT=no-change
SECOND_RUN_RESULT=changed
FINAL_VALUE=release-2026.10.2
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac004EndToEnd proves the Practice lab's starting state fails (the
// resource is still dev-shaped) and the documented solution — running
// infra-apply against the prod environment — passes.
func TestIac004EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-004", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-004"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("resource still dev-shaped: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
infra-apply prod
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac001EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — editing
// resources.yaml's declared content and re-running infra-apply instead
// of editing /etc/banner.txt directly — passes.
func TestIac001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
before=$(sed -n 's/^content: *//p' /infra/resources.yaml)

sed -i 's/^content: .*/content: Build server banner v2/' /infra/resources.yaml

infra-apply

after=$(sed -n 's/^content: *//p' /infra/resources.yaml)

cat <<EOF > /workspace/answers.env
DESIRED_VALUE_BEFORE="$before"
DESIRED_VALUE_AFTER="$after"
SOURCE_OF_TRUTH_FILE=/infra/resources.yaml
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac002EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet, and resources.yaml already differs from the real
// resource) and the documented solution — recording infra-apply
// --plan's output before applying for real — passes.
func TestIac002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
before=$(cat /infra/managed-file.txt)

plan_output=$(infra-apply --plan)

after_plan=$(cat /infra/managed-file.txt)
if [ "$before" != "$after_plan" ]; then
  echo "plan should not have changed the real resource" >&2
  exit 1
fi

cat <<EOF > /workspace/answers.env
PLANNED_CHANGE="$plan_output"
EOF

infra-apply
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac005EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet, and infra-apply refuses to run under a stale lock)
// and the documented solution — confirming the locked PID is dead,
// removing the stale lock, and re-running infra-apply — passes.
func TestIac005EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-005", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-005"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
lock_pid=$(sed -n 's/^PID=//p' /var/run/infra/apply.lock)

if ps -p "$lock_pid" >/dev/null 2>&1; then
  echo "unexpected: lock PID is alive" >&2
  exit 1
fi

rm /var/run/infra/apply.lock

infra-apply

cat <<EOF > /workspace/answers.env
LOCK_PID=$lock_pid
LOCK_PROCESS_ALIVE=no
APPLY_SUCCEEDED=yes
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac006EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — recording the file
// and environment values, then running infra-apply and reading which
// source it actually resolved to — passes.
func TestIac006EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-006", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-006"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
file_value=$(sed -n 's/^value: *//p' /infra/resources.yaml)
env_value="$INFRA_VALUE"

infra-apply
infra-apply --value=maintenance-mode-forced

cat <<EOF > /workspace/answers.env
FILE_VALUE=$file_value
ENV_VALUE=$env_value
ACTUAL_SOURCE=env
RESOLVED_VALUE=$env_value
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac007EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet, and the pin disagrees with the declared module) and
// the documented solution — fixing the pin to match resources.yaml's
// declared module and re-running infra-apply — passes.
func TestIac007EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-007", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-007"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
declared=$(sed -n 's/^module: *//p' /infra/resources.yaml)
pinned=$(cat /infra/MODULE_VERSION)

echo "$declared" > /infra/MODULE_VERSION

infra-apply

cat <<EOF > /workspace/answers.env
DECLARED_MODULE=$declared
PINNED_MODULE=$pinned
MODULES_AGREED_BEFORE_FIX=no
RESOLVED_VERSION_AFTER_FIX=$declared
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac008EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet, and the cache resource is skipped as unmanaged) and
// the documented solution — importing the missing path into state and
// re-running infra-apply — passes.
func TestIac008EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-008", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-008"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
infra-apply --import /etc/cache/cache.conf
infra-apply

cat <<EOF > /workspace/answers.env
UNMANAGED_PATH=/etc/cache/cache.conf
APPLY_SUCCEEDED_AFTER_IMPORT=yes
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs001EndToEnd proves the Guided lab's starting state fails (no
// answers.env yet) and the documented solution — using jq to tell a
// discrete log event apart from a current metric value — passes.
func TestObs001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
ERROR_COUNT=3
REQUEST_ID_MESSAGE="upstream inventory service timed out"
DISK_USED_PCT=78
WHY_ANSWER=LOG
HOWMUCH_ANSWER=METRIC
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs002EndToEnd proves the Practice lab's starting state fails (no
// answers.env yet) and the documented solution — extracting the same
// fact independently from an unstructured log and a structured one —
// passes.
func TestObs002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("no answers.env yet: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<EOF > /workspace/answers.env
ATTEMPTS_FROM_LEGACY=3
FINAL_FAILURE_TIME_FROM_LEGACY="2026-01-15 10:22:07"
METHOD_LEGACY=grep+awk
ATTEMPTS_FROM_JSON=3
FINAL_FAILURE_TIME_FROM_JSON="2026-01-15 10:22:07"
METHOD_JSON=jq
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac009EndToEnd proves the Practice lab's starting state fails
// (main.tf was edited but never applied) and the documented solution
// passes.
func TestIac009EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-009", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-009"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("drifted banner: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /infra
terraform apply -auto-approve -input=false
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac010EndToEnd proves the Practice lab's starting state fails
// (the renamed resource was never reconciled in state) and the
// documented solution passes.
func TestIac010EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-010", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-010"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("unreconciled rename: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cd /infra
terraform state mv local_file.legal_notice local_file.privacy_notice
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestIac011EndToEnd proves the Practice lab's starting state fails
// (the shell task is never idempotent) and the documented solution
// passes.
func TestIac011EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "iac-011", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-iac-011"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("non-idempotent playbook: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /infra/playbook.yml
---
- hosts: localhost
  connection: local
  tasks:
    - name: ensure max_connections is configured
      ansible.builtin.lineinfile:
        path: /etc/app/app.conf
        regexp: '^max_connections='
        line: 'max_connections=200'
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd008EndToEnd proves the break-fix lab's starting state fails
// (deploy's output is declared but never set) and the documented
// solution passes.
func TestCicd008EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-008", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-008"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("deploy silently skipped: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /workspace/project/.github/workflows/ci.yml
name: CI
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      ready: ${{ steps.build.outputs.ready }}
    steps:
      - name: Build the app
        id: build
        run: |
          mkdir -p /workspace/dist
          echo "built" > /workspace/dist/app.txt
          echo "ready=true" >> "$GITHUB_OUTPUT"

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: ${{ needs.build.outputs.ready == 'true' }}
    steps:
      - name: Deploy
        run: |
          cp /workspace/dist/app.txt /workspace/deployed.txt
          echo "deployed"
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestCicd009EndToEnd proves the break-fix lab's starting state fails
// (deploy's rules: only matches a tag pipeline) and the documented
// solution passes.
func TestCicd009EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "cicd-009", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-cicd-009"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("deploy silently skipped: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /workspace/project/.gitlab-ci.yml
stages:
  - build
  - deploy

build:
  stage: build
  script:
    - mkdir -p /workspace/dist
    - echo "built" > /workspace/dist/app.txt

deploy:
  stage: deploy
  needs: ["build"]
  rules:
    - if: '$CI_COMMIT_BRANCH'
  script:
    - cp /workspace/dist/app.txt /workspace/deployed.txt
    - echo "deployed"
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs008EndToEnd proves the break-fix lab's starting state fails
// (the real PromQL expr checks the wrong direction) and the documented
// solution passes.
func TestObs008EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-008", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-008"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("inverted alert expr: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
sed -i 's/error_rate < 5/error_rate > 5/' /infra/rules.yml
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestObs009EndToEnd proves the break-fix lab's starting state fails
// (the dashboard panel queries a renamed metric) and the documented
// solution passes, once Grafana's file provisioning poller has picked
// up the change.
func TestObs009EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "obs-009", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-obs-009"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("stale panel query: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
cat <<'EOF' > /etc/grafana/provisioning/dashboards/json/checkout.json
{
  "uid": "checkout-overview",
  "title": "Checkout overview",
  "panels": [
    {
      "id": 1,
      "title": "HTTP requests",
      "type": "timeseries",
      "gridPos": { "h": 8, "w": 12, "x": 0, "y": 0 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "http_requests_count_total",
          "refId": "A",
          "datasource": { "type": "prometheus", "uid": "prometheus" }
        }
      ]
    }
  ],
  "schemaVersion": 39,
  "version": 2
}
EOF
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	// Grafana polls its provisioning directory on an interval rather
	// than reloading instantly -- give it time to pick up the edit
	// before asserting the fix landed.
	time.Sleep(15 * time.Second)

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestAzure001EndToEnd proves the break-fix lab's starting state fails
// (retentionDays' default value is a string, not an int) and the
// documented solution passes.
func TestAzure001EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "azure-001", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-azure-001"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("bicep type error: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
sed -i "s/param retentionDays int = '30'/param retentionDays int = 30/" /infra/main.bicep
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}

// TestAzure002EndToEnd proves the break-fix lab's starting state fails
// (the container's public access level blocks anonymous reads) and the
// documented solution passes.
func TestAzure002EndToEnd(t *testing.T) {
	if os.Getenv("BREAKLAB_INTEGRATION") != "1" {
		t.Skip("set BREAKLAB_INTEGRATION=1 to run against real Docker")
	}
	ctx := context.Background()
	r := newIntegrationRunner(t)
	if err := r.Start(ctx, "azure-002", true); err != nil {
		t.Fatal(err)
	}
	container := "breaklab-azure-002"

	if res, err := r.Check(ctx); err != nil || res.Passed {
		t.Fatalf("container not publicly readable: check should fail, got %+v, %v", res, err)
	}

	solution := `set -e
az storage container set-permission --name public-assets --public-access blob
`
	if out, err := r.Docker.Exec(ctx, container, "bash", "-c", solution); err != nil {
		t.Fatalf("applying solution: %v\n%s", err, out)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("fixed lab should pass: %+v, %v", res, err)
	}
}
