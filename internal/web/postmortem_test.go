package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"breaklab/internal/web/postmortems"
	"breaklab/internal/web/webtest"
)

func TestPostmortemRouteReturns404ForUnknownOrTraversalSlug(t *testing.T) {
	dir := t.TempDir()
	webtest.MustWriteFile(t, filepath.Join(dir, "progress", "progress.json"), `{"version":1,"attempts":[]}`)
	templatesDir := t.TempDir()
	cfg := Config{
		ProgressPath: filepath.Join(dir, "progress", "progress.json"),
		TemplatesDir: templatesDir,
	}

	for _, slug := range []string{"no-such-postmortem", "..", "etc%2Fpasswd"} {
		req := httptest.NewRequest("GET", "/postmortems/"+slug, nil)
		rec := httptest.NewRecorder()
		Routes(cfg).ServeHTTP(rec, req)
		switch rec.Code {
		case http.StatusNotFound, http.StatusMovedPermanently, http.StatusTemporaryRedirect:
			// Either this package's own 404 (unknown slug) or net/http's
			// own path-cleaning redirect for a ".."-shaped segment —
			// both mean no file was ever read for it.
		default:
			t.Errorf("slug %q: status = %d, want 404 or a path-cleaning redirect", slug, rec.Code)
		}
	}
}

// TestPostmortemsIndexAPIServesFilteredPage is the `frontend`-capability
// counterpart of the old TestPostmortemsIndexHandlerServesFilteredPage:
// "/postmortems" is a static Astro shell now, so this checks the JSON
// API's data instead of rendered markup (and the sidebar, which that
// API never sets — see postmortemsIndexAPIHandler — isn't checked here
// at all).
func TestPostmortemsIndexAPIServesFilteredPage(t *testing.T) {
	cfg := fullFixtureConfig(t)

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/postmortems", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var view postmortems.PostmortemsIndexView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, row := range view.Rows {
		if row.Lab == "The website returns nothing" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected the fixture lab's postmortem row, got: %+v", view.Rows)
	}

	rec = httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/postmortems?stage=git", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	view = postmortems.PostmortemsIndexView{}
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	for _, row := range view.Rows {
		if row.Lab == "The website returns nothing" {
			t.Error("stage=git filter should have excluded the linux-stage lab")
		}
	}
}
