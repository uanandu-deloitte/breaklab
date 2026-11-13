package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDashboardAPINotMountedWithoutFrontendDir(t *testing.T) {
	cfg := testConfig(t, 0) // FrontendDir left empty

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/dashboard", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404 (no FrontendDir configured)", rec.Code)
	}
}
