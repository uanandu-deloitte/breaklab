package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"breaklab/internal/progress"
)

func TestSyncOnceRoundTrip(t *testing.T) {
	srv, store := newFakeAccountServer(t, "good-token", nil)
	link := &accountLink{Server: srv.URL, Token: "good-token"}

	pathA := filepath.Join(t.TempDir(), "progress.json")
	solved := time.Now().UTC()
	if err := progress.Save(pathA, &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{ID: "a1", Lab: "web-down", Status: progress.StatusSolved, StartedAt: solved.Add(-time.Minute), SolvedAt: &solved},
	}}); err != nil {
		t.Fatal(err)
	}

	pushed, pulled, err := syncOnce(pathA, link)
	if err != nil {
		t.Fatal(err)
	}
	if pushed != 1 || pulled != 0 {
		t.Fatalf("machine A's first sync: pushed=%d pulled=%d, want 1, 0", pushed, pulled)
	}
	if len(*store) != 1 {
		t.Fatalf("hosted store = %+v, want exactly the one pushed attempt", *store)
	}

	pathB := filepath.Join(t.TempDir(), "progress.json")
	pushedB, pulledB, err := syncOnce(pathB, link)
	if err != nil {
		t.Fatal(err)
	}
	if pushedB != 0 || pulledB != 1 {
		t.Fatalf("machine B's first sync (fresh progress.json): pushed=%d pulled=%d, want 0, 1", pushedB, pulledB)
	}
	pfB, err := progress.Load(pathB)
	if err != nil {
		t.Fatal(err)
	}
	if len(pfB.Attempts) != 1 || pfB.Attempts[0].ID != "a1" {
		t.Fatalf("machine B's progress.json = %+v, want a1 pulled down", pfB.Attempts)
	}

	// Re-syncing either machine must never duplicate anything, locally
	// or on the hosted side.
	if _, pulledAgain, err := syncOnce(pathA, link); err != nil || pulledAgain != 0 {
		t.Fatalf("re-syncing machine A: pulled=%d, err=%v, want 0 pulled (a1 is already local)", pulledAgain, err)
	}
	if len(*store) != 1 {
		t.Fatalf("hosted store after re-sync = %+v, want still exactly one attempt", *store)
	}
	pfBAgain, err := progress.Load(pathB)
	if err != nil {
		t.Fatal(err)
	}
	if len(pfBAgain.Attempts) != 1 {
		t.Fatalf("machine B's progress.json after a second sync = %+v, want still exactly one attempt", pfBAgain.Attempts)
	}
}

func TestSyncNeverPushesActiveAttempts(t *testing.T) {
	srv, store := newFakeAccountServer(t, "good-token", nil)
	link := &accountLink{Server: srv.URL, Token: "good-token"}

	path := filepath.Join(t.TempDir(), "progress.json")
	if err := progress.Save(path, &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{ID: "active1", Lab: "web-down", Status: progress.StatusActive, StartedAt: time.Now().UTC()},
	}}); err != nil {
		t.Fatal(err)
	}

	pushed, _, err := syncOnce(path, link)
	if err != nil {
		t.Fatal(err)
	}
	if pushed != 0 {
		t.Errorf("pushed = %d, want 0: an active attempt must never cross the wire", pushed)
	}
	if len(*store) != 0 {
		t.Errorf("hosted store = %+v, want empty", *store)
	}
}

func TestRunSyncNoTokenLeavesProgressUntouched(t *testing.T) {
	withAccountLinkPath(t, filepath.Join(t.TempDir(), ".breaklab", "account.json"))
	progPath := filepath.Join(t.TempDir(), "progress.json")
	withSyncProgressPath(t, progPath)

	before := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{ID: "x", Lab: "web-down", Status: progress.StatusSolved},
	}}
	if err := progress.Save(progPath, before); err != nil {
		t.Fatal(err)
	}

	if code := runSync(nil); code == 0 {
		t.Fatal("runSync with no linked account should exit non-zero")
	}

	after, err := progress.Load(progPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Attempts) != 1 || after.Attempts[0].ID != "x" {
		t.Fatalf("progress.json changed while not logged in: %+v", after.Attempts)
	}
}

// TestSyncNeverImportsDockerOrRunner is the structural half of
// verifying `breaklab sync` never starts, stops, or queries Docker
// (design.md "Sync only ever moves terminal attempts" / sync spec
// "Sync never touches lab execution"): sync.go doesn't even import
// the packages that could do that, so there is nothing to fake out.
// The functional round-trip test above is the behavioral half — it
// runs with no Docker daemon available at all and still passes.
func TestSyncNeverImportsDockerOrRunner(t *testing.T) {
	data, err := os.ReadFile("sync.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"breaklab/internal/docker", "breaklab/internal/runner"} {
		if strings.Contains(string(data), forbidden) {
			t.Errorf("sync.go must never import %s", forbidden)
		}
	}
}
