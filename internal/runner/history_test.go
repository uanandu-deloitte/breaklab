package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHistoryRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".breaklab", "history.json")

	if v, err := LastVariant(path, "web-down"); v != "" || err != nil {
		t.Fatalf("missing file: v = %q, err = %v", v, err)
	}

	if err := RecordVariant(path, "web-down", "config-typo"); err != nil {
		t.Fatal(err)
	}
	if v, err := LastVariant(path, "web-down"); err != nil || v != "config-typo" {
		t.Fatalf("v = %q, err = %v", v, err)
	}

	// A different lab's entry is preserved alongside, and overwriting
	// web-down's entry doesn't disturb it.
	if err := RecordVariant(path, "disk-full", "runaway-log"); err != nil {
		t.Fatal(err)
	}
	if err := RecordVariant(path, "web-down", "port-taken"); err != nil {
		t.Fatal(err)
	}
	if v, err := LastVariant(path, "web-down"); err != nil || v != "port-taken" {
		t.Fatalf("v = %q, err = %v", v, err)
	}
	if v, err := LastVariant(path, "disk-full"); err != nil || v != "runaway-log" {
		t.Fatalf("v = %q, err = %v", v, err)
	}

	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, err %v, want 0600", info, err)
	}
}

func TestHistoryCorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LastVariant(path, "web-down"); err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("err = %v", err)
	}

	// RecordVariant tolerates the corrupt file by discarding it.
	if err := RecordVariant(path, "web-down", "config-typo"); err != nil {
		t.Fatal(err)
	}
	if v, err := LastVariant(path, "web-down"); err != nil || v != "config-typo" {
		t.Fatalf("v = %q, err = %v", v, err)
	}
}
