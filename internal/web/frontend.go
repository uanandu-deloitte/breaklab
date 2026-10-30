package web

import (
	"net/http"
	"path/filepath"
)

// frontendShellHandler serves shellPath (one static HTML file the
// Astro build produced from a single placeholder parameter value) for
// every request to a path-parameterized route, regardless of the
// actual parameter value in the URL — design.md's "Astro builds
// shells, not content" decision. The real id/slug/name is resolved
// client-side from the URL and is never needed to produce this
// response.
func frontendShellHandler(cfg Config, shellPath string) http.HandlerFunc {
	path := filepath.Join(cfg.FrontendDir, shellPath)
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

// frontendAssetHandler serves the Astro build's hashed JS/CSS chunks
// (frontend/dist/_astro/*) under /_astro/ — the one directory every
// Astro static build writes its assets to, shared by every migrated
// page regardless of which one requested them.
func frontendAssetHandler(cfg Config) http.Handler {
	dir := filepath.Join(cfg.FrontendDir, "_astro")
	return http.StripPrefix("/_astro/", http.FileServer(http.Dir(dir)))
}
