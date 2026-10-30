package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRawSession(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLoadSessionRejectsTamperedFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")
	cases := map[string]string{
		"flag-shaped container": `{"lab":"alpha","container":"-f","variant":"v"}`,
		"other container":       `{"lab":"alpha","container":"breaklab-other","variant":"v"}`,
		"traversal lab":         `{"lab":"../x","container":"breaklab-../x","variant":"v"}`,
		"traversal variant":     `{"lab":"alpha","container":"breaklab-alpha","variant":"../../x"}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			writeRawSession(t, path, body)
			if _, err := LoadSession(path); err == nil || !strings.Contains(err.Error(), path) {
				t.Errorf("err = %v, want a corrupt-session error naming %s", err, path)
			}
		})
	}
}

// TestLoadSessionAcceptsEmptyVariant confirms a session for a
// no-variant lab (empty Variant, see selectVariant) loads normally —
// catalog.ValidID would reject an empty string, so LoadSession must
// special-case it rather than treating every no-variant session as
// corrupt.
func TestLoadSessionAcceptsEmptyVariant(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")
	writeRawSession(t, path, `{"lab":"solo","container":"breaklab-solo","variant":""}`)
	s, err := LoadSession(path)
	if err != nil || s == nil || s.Variant != "" {
		t.Fatalf("session = %+v, err = %v", s, err)
	}
}

func TestStartWrittenSessionStillLoads(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	if err := r.Start(context.Background(), "alpha", false); err != nil {
		t.Fatal(err)
	}
	if s, err := LoadSession(r.SessionPath); err != nil || s.Container != "breaklab-alpha" {
		t.Fatalf("session = %+v, err = %v", s, err)
	}
}

func TestStopOnTamperedSessionTouchesNoDocker(t *testing.T) {
	r, fd, _ := setup(t, "alpha")
	writeRawSession(t, r.SessionPath, `{"lab":"alpha","container":"-f","variant":"first-variant"}`)

	if err := r.Stop(context.Background()); err == nil {
		t.Fatal("expected Stop to fail on a tampered session")
	}
	if len(fd.calls) != 0 {
		t.Errorf("docker was called: %v", fd.calls)
	}
}
