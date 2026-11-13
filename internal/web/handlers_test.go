package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthzReturns200 covers the Kubernetes liveness/readiness probe
// target: it must answer 200 regardless of whether a frontend build is
// configured, since that's not a dependency it needs.
func TestHealthzReturns200(t *testing.T) {
	for name, cfg := range map[string]Config{
		"zero value":      {},
		"fully populated": actionTestConfig(t, nil),
	} {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/healthz", nil))
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
			}
		})
	}
}
