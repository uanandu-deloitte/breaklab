package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRootServesAstroShellWhenFrontendDirSet(t *testing.T) {
	cfg := actionTestConfig(t, nil)
	cfg.FrontendDir = t.TempDir()
	if err := os.MkdirAll(filepath.Join(cfg.FrontendDir, "dashboard"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.FrontendDir, "dashboard", "index.html"), []byte("<html>astro dashboard</html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "<html>astro dashboard</html>" {
		t.Fatalf("code = %d, body = %q, want the Astro shell's content, not the old template render", rec.Code, rec.Body.String())
	}
}

// TestRootWithNoFrontendDistIs404 pins the deliberate consequence of
// deleting the dashboard's html/template fallback (redesign-web-ui-astro
// section 3.4): without a built frontend/dist, "/" has nothing to
// serve. There is no dual rendering path to fall back to anymore.
func TestRootWithNoFrontendDistIs404(t *testing.T) {
	cfg := actionTestConfig(t, nil) // FrontendDir left empty, the default

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404 (no FrontendDir, and no template fallback exists anymore)", rec.Code)
	}
}
