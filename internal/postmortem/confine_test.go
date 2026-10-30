package postmortem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfined(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "postmortems")
	write := func(path string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	good := filepath.Join(dir, "2026-10-17-web-down-config-typo.md")
	write(good)
	write(filepath.Join(root, "outside.md"))
	write(filepath.Join(dir, "sub", "nested.md"))
	write(filepath.Join(dir, "notes.txt"))
	if err := os.Symlink("/etc/passwd", filepath.Join(dir, "link.md")); err != nil {
		t.Fatal(err)
	}

	if got, err := Confined(dir, good); err != nil || got != good {
		t.Errorf("absolute direct child: got %q, %v", got, err)
	}

	// A relative recorded path (as the CLI writes them) is accepted too.
	t.Chdir(root)
	if _, err := Confined("postmortems", filepath.Join("postmortems", filepath.Base(good))); err != nil {
		t.Errorf("relative direct child: %v", err)
	}

	for _, bad := range []string{
		filepath.Join(dir, "..", "outside.md"),
		filepath.Join(dir, "sub", "nested.md"),
		"/etc/passwd",
		filepath.Join(dir, "notes.txt"),
		filepath.Join(dir, "link.md"),
		filepath.Join(dir, "missing.md"),
	} {
		if _, err := Confined(dir, bad); err == nil {
			t.Errorf("Confined(%q) accepted, want rejected", bad)
		}
	}
}
