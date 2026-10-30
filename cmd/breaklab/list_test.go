package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Task 4.2: a broken second lab produces one warning on stderr but does not
// stop the valid lab from being listed.
func TestListWarnsOnBrokenLabWithoutBlockingOthers(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "breaklab")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building breaklab: %v\n%s", err, out)
	}

	workdir := t.TempDir()
	writeFile(t, filepath.Join(workdir, "roadmap.yaml"), `
stages:
  - id: linux
    title: Linux fundamentals
`)
	writeFile(t, filepath.Join(workdir, "labs", "web-down", "lab.yaml"), `
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
`)
	writeFile(t, filepath.Join(workdir, "labs", "broken-lab", "lab.yaml"), `
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
`)

	cmd := exec.Command(bin, "list")
	cmd.Dir = workdir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("breaklab list failed: %v (stderr: %s)", err, stderr.String())
	}

	if !contains(stdout.String(), "web-down") {
		t.Fatalf("expected web-down in stdout, got: %s", stdout.String())
	}
	if !contains(stderr.String(), "broken-lab") {
		t.Fatalf("expected a warning naming broken-lab on stderr, got: %s", stderr.String())
	}
}

// Task 2.1 (add-local-web-ui): a lab with a complete attempt in
// progress.json shows status "complete" instead of the always-
// "available" stub from change 1.
func TestListShowsAccurateStatusFromProgress(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "breaklab")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building breaklab: %v\n%s", err, out)
	}

	workdir := t.TempDir()
	writeFile(t, filepath.Join(workdir, "roadmap.yaml"), `
stages:
  - id: linux
    title: Linux fundamentals
`)
	writeFile(t, filepath.Join(workdir, "labs", "web-down", "lab.yaml"), `
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
`)
	writeFile(t, filepath.Join(workdir, "progress", "progress.json"), `
{
  "version": 1,
  "attempts": [
    {
      "lab": "web-down",
      "variant": "config-typo",
      "started_at": "2026-10-17T10:02:00Z",
      "solved_at": "2026-10-17T10:20:00Z",
      "failed_checks": 1,
      "hints_used": 1,
      "postmortem": "postmortems/2026-10-17-web-down-config-typo.md",
      "status": "complete"
    }
  ]
}
`)

	cmd := exec.Command(bin, "list")
	cmd.Dir = workdir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("breaklab list failed: %v (stderr: %s)", err, stderr.String())
	}
	if !contains(stdout.String(), "complete") {
		t.Fatalf("expected status \"complete\" in stdout, got: %s", stdout.String())
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
