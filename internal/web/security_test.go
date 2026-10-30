package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"breaklab/internal/web/webtest"
)

// A hand-edited progress.json pointing a complete attempt's postmortem
// outside postmortems/ must be neither served nor linked.
func TestPostmortemOutsideDirIsNotServedOrLinked(t *testing.T) {
	cfg := fullFixtureConfig(t)
	root := filepath.Dir(cfg.LabsDir)

	outside := filepath.Join(root, "elsewhere", "2026-10-17-web-down-config-typo.md")
	webtest.MustWriteFile(t, outside, strings.Replace(postmortemFixtureBody, "homepage stopped responding", "OUTSIDE_FILE_MARKER", 1))
	if err := os.Remove(filepath.Join(cfg.PostmortemsDir, "2026-10-17-web-down-config-typo.md")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(cfg.ProgressPath)
	if err != nil {
		t.Fatal(err)
	}
	inside := filepath.Join(cfg.PostmortemsDir, "2026-10-17-web-down-config-typo.md")
	webtest.MustWriteFile(t, cfg.ProgressPath, strings.Replace(string(data), inside, outside, 1))

	get := func(path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		return rec
	}

	// "/postmortems/{slug}" is a static Astro shell now (no embedded
	// data) — the real 404-for-an-invalid-slug check lives at the JSON
	// API layer instead.
	if rec := get("/api/postmortems/2026-10-17-web-down-config-typo"); rec.Code != http.StatusNotFound {
		t.Errorf("slug API status = %d, want 404", rec.Code)
	} else if strings.Contains(rec.Body.String(), "OUTSIDE_FILE_MARKER") {
		t.Error("out-of-dir postmortem content was served")
	}
	for _, page := range []string{"/", "/roadmap", "/postmortems"} {
		rec := get(page)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d", page, rec.Code)
		}
		if strings.Contains(rec.Body.String(), "/postmortems/") {
			t.Errorf("%s links to a postmortem whose recorded path is outside postmortems/", page)
		}
	}
}

// The web UI serves exactly these read-only routes; nothing else maps to
// a file, and no route acts on a path smuggled in through the query.
func TestRouteTableIsReadOnlyAndBounded(t *testing.T) {
	cfg := fullFixtureConfig(t)
	slug := "2026-10-17-web-down-config-typo"
	served := []string{"/", "/roadmap", "/challenges", "/labs/web-down", "/postmortems", "/postmortems/" + slug, "/js/cytoscape.min.js", "/topology"}

	passwd, _ := os.ReadFile("/etc/passwd")
	marker := ""
	if len(passwd) > 0 {
		marker = strings.SplitN(string(passwd), "\n", 2)[0]
	}

	for _, p := range served {
		for _, q := range []string{"", "?path=/etc/passwd", "?file=../../etc/passwd&cmd=id"} {
			rec := httptest.NewRecorder()
			Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", p+q, nil))
			if rec.Code != http.StatusOK {
				t.Errorf("GET %s%s: status %d, want 200", p, q, rec.Code)
			}
			if marker != "" && strings.Contains(rec.Body.String(), marker) {
				t.Errorf("GET %s%s: response contains /etc/passwd content", p, q)
			}
		}
		for _, method := range []string{"POST", "PUT", "DELETE"} {
			rec := httptest.NewRecorder()
			Routes(cfg).ServeHTTP(rec, httptest.NewRequest(method, p, nil))
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("%s %s: status %d, want 405", method, p, rec.Code)
			}
		}
	}

	// "/postmortems/<slug>.md" is deliberately absent from this list:
	// "/postmortems/{slug}" is a static Astro shell now
	// (frontendShellHandler) that serves the exact same fixed file for
	// every path-parameter value, .md-suffixed or not — there is no
	// code path by which any slug value reads a different file, so the
	// old "appending .md shouldn't serve the raw file" check has no
	// remaining way to fail, by construction rather than by a 404.
	for _, p := range []string{
		"/labs/web-down/TICKET.md", "/labs/web-down/variants/config-typo.sh",
		"/solutions/config-typo.sh", "/variants", "/progress/progress.json",
		"/templates/web/lab.html", "/roadmap.yaml",
	} {
		rec := httptest.NewRecorder()
		Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", p, nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("GET %s: status %d, want 404", p, rec.Code)
		}
	}
}
