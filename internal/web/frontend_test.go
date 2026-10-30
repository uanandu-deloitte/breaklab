package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFrontendAssetServedWhenFrontendDirSet(t *testing.T) {
	cfg := testConfig(t, 0)
	cfg.FrontendDir = t.TempDir()
	if err := os.MkdirAll(filepath.Join(cfg.FrontendDir, "_astro"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.FrontendDir, "_astro", "chunk.js"), []byte("console.log(1)"), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/_astro/chunk.js", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "console.log(1)" {
		t.Fatalf("code = %d, body = %q", rec.Code, rec.Body.String())
	}
}

func TestFrontendShellHandlerIgnoresPathParameter(t *testing.T) {
	cfg := testConfig(t, 0)
	cfg.FrontendDir = t.TempDir()
	if err := os.MkdirAll(filepath.Join(cfg.FrontendDir, "things", "_"), 0o755); err != nil {
		t.Fatal(err)
	}
	shellPath := filepath.Join("things", "_", "index.html")
	if err := os.WriteFile(filepath.Join(cfg.FrontendDir, shellPath), []byte("<html>shell</html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /things/{id}", frontendShellHandler(cfg, shellPath))

	recA := httptest.NewRecorder()
	mux.ServeHTTP(recA, httptest.NewRequest("GET", "/things/abc", nil))
	recB := httptest.NewRecorder()
	mux.ServeHTTP(recB, httptest.NewRequest("GET", "/things/some-other-id", nil))

	if recA.Code != http.StatusOK || recB.Code != http.StatusOK {
		t.Fatalf("codes = %d, %d, want both 200", recA.Code, recB.Code)
	}
	if recA.Body.String() != recB.Body.String() {
		t.Fatalf("bodies differ: %q vs %q, want byte-identical (the real id is never needed server-side)", recA.Body.String(), recB.Body.String())
	}
	if recA.Body.String() != "<html>shell</html>" {
		t.Fatalf("body = %q, want the shell file's content", recA.Body.String())
	}
}

func TestFrontendRouteNotMountedWithoutFrontendDir(t *testing.T) {
	cfg := testConfig(t, 0) // FrontendDir left empty, the default

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/_astro/chunk.js", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404 (no FrontendDir configured)", rec.Code)
	}

	// A vendored static asset (served from TemplatesDir, unrelated to
	// FrontendDir) is untouched by FrontendDir being unset.
	if err := os.MkdirAll(filepath.Join(cfg.TemplatesDir, "js"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.TemplatesDir, "js", "cytoscape.min.js"), []byte("cytoscape-stub"), 0o644); err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/js/cytoscape.min.js", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("/js/cytoscape.min.js code = %d, want 200 (unaffected by this change)", rec.Code)
	}
}
