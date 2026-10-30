// Package webtest holds tiny test-fixture helpers shared by more than
// one internal/web domain package's tests — not production code, but
// Go has no other way to share an unexported helper across package
// boundaries.
package webtest

import (
	"os"
	"path/filepath"
	"testing"
)

// MustWriteFile writes content to path, creating any missing parent
// directories, failing the test on any error.
func MustWriteFile(t testing.TB, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
