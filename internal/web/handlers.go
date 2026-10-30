package web

import (
	"net/http"
	"path/filepath"

	"breaklab/internal/web/account"
	"breaklab/internal/web/actions"
	"breaklab/internal/web/capstone"
	"breaklab/internal/web/challenges"
	"breaklab/internal/web/chaos"
	"breaklab/internal/web/dashboard"
	"breaklab/internal/web/lab"
	"breaklab/internal/web/postmortems"
	"breaklab/internal/web/roadmap"
	"breaklab/internal/web/topology"
)

// Routes builds the server's handler. Each route is registered
// explicitly; there is no generic static file server over any
// directory, so a request can never reach anything this package doesn't
// already know how to serve.
func Routes(cfg Config) http.Handler {
	mux := http.NewServeMux()
	if cfg.FrontendDir != "" {
		mux.Handle("GET /_astro/", frontendAssetHandler(cfg))
		dashboard.RegisterRoutes(mux, cfg)
		roadmap.RegisterRoutes(mux, cfg)
		challenges.RegisterRoutes(mux, cfg)
		lab.RegisterRoutes(mux, cfg)
		actions.RegisterRoutes(mux, cfg)
		topology.RegisterRoutes(mux, cfg)
		postmortems.RegisterRoutes(mux, cfg)
		chaos.RegisterRoutes(mux, cfg)
		capstone.RegisterRoutes(mux, cfg)
	}
	if cfg.FrontendDir != "" && cfg.AccountStore != nil {
		account.RegisterRoutes(mux, cfg)
	}
	mux.HandleFunc("GET /fonts/JetBrainsMono-Regular.woff2", staticAssetHandler(cfg, "fonts/JetBrainsMono-Regular.woff2", "font/woff2"))
	mux.HandleFunc("GET /fonts/JetBrainsMono-Medium.woff2", staticAssetHandler(cfg, "fonts/JetBrainsMono-Medium.woff2", "font/woff2"))
	mux.HandleFunc("GET /fonts/Geist-Regular.woff2", staticAssetHandler(cfg, "fonts/Geist-Regular.woff2", "font/woff2"))
	mux.HandleFunc("GET /fonts/Geist-SemiBold.woff2", staticAssetHandler(cfg, "fonts/Geist-SemiBold.woff2", "font/woff2"))
	mux.HandleFunc("GET /fonts/Geist-Bold.woff2", staticAssetHandler(cfg, "fonts/Geist-Bold.woff2", "font/woff2"))
	mux.HandleFunc("GET /js/cytoscape.min.js", staticAssetHandler(cfg, "js/cytoscape.min.js", "text/javascript; charset=utf-8"))
	// Migrated to the `frontend` capability (Astro); no html/template
	// fallback — same as "/" (see the comment on that route below).
	mux.HandleFunc("GET /roadmap", frontendShellHandler(cfg, filepath.Join("roadmap", "index.html")))
	mux.HandleFunc("GET /chaos", frontendShellHandler(cfg, filepath.Join("chaos", "index.html")))
	mux.HandleFunc("GET /capstone", frontendShellHandler(cfg, filepath.Join("capstone", "index.html")))
	mux.HandleFunc("GET /challenges", frontendShellHandler(cfg, filepath.Join("challenges", "index.html")))
	mux.HandleFunc("GET /labs/{id}", frontendShellHandler(cfg, filepath.Join("labs", "_", "index.html")))
	// Migrated to the `frontend` capability (Astro); no html/template
	// fallback.
	mux.HandleFunc("GET /topology", frontendShellHandler(cfg, filepath.Join("topology", "index.html")))
	mux.HandleFunc("GET /postmortems", frontendShellHandler(cfg, filepath.Join("postmortems", "index.html")))
	mux.HandleFunc("GET /postmortems/{slug}", frontendShellHandler(cfg, filepath.Join("postmortems", "_", "index.html")))
	// The dashboard migrated to the `frontend` capability (Astro); there
	// is no html/template fallback — breaklab serve requires frontend/dist
	// to exist (see --frontend's default and docs/deploy-raspberry-pi.md).
	mux.HandleFunc("GET /{$}", frontendShellHandler(cfg, filepath.Join("dashboard", "index.html"))) // exact root only; unknown paths 404

	// Hosted-account routes only exist when a Store was configured
	// (`breaklab serve --db`); with no AccountStore, none of this is
	// even registered, so local-only usage is completely unaffected.
	// The GET pages are static Astro shells (frontend capability), same
	// as every other page registered above.
	if cfg.AccountStore != nil {
		mux.HandleFunc("GET /signup", frontendShellHandler(cfg, filepath.Join("signup", "index.html")))
		mux.HandleFunc("GET /login", frontendShellHandler(cfg, filepath.Join("login", "index.html")))
		mux.HandleFunc("GET /account", frontendShellHandler(cfg, filepath.Join("account", "index.html")))
		mux.HandleFunc("GET /u/{name}", frontendShellHandler(cfg, filepath.Join("u", "_", "index.html")))
	}
	return mux
}

// staticCacheControl is applied to every fixed, vendored asset (a font
// or vendored JS file) — none of them change except on a redeploy
// (which restarts the process and gets a fresh browser session anyway
// in practice), so a short private cache is a free repeat-visit win
// without risking a stale asset surviving a real update for long.
// http.ServeFile still handles conditional requests (ETag/If-Modified-
// Since) on top of this for anything that revalidates.
const staticCacheControl = "private, max-age=3600"

// staticAssetHandler serves one fixed, vendored static asset (a font or
// vendored JS file) under cfg.TemplatesDir — always one path baked in
// at registration time, never anything derived from the request.
func staticAssetHandler(cfg Config, relPath, contentType string) http.HandlerFunc {
	path := filepath.Join(cfg.TemplatesDir, filepath.FromSlash(relPath))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", staticCacheControl)
		http.ServeFile(w, r, path)
	}
}
