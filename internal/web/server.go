// Package web serves BreakLab's local dashboard/roadmap/challenge
// browser/lab detail/postmortem viewer: the view a learner opens every
// day, reading progress.json, roadmap.yaml, lab metadata, and
// postmortems live, with no publish step and no publish-flag gating.
//
// This package is the composition root: it owns Serve/Routes and
// nothing else — every page's handlers live in their own domain
// subpackage (dashboard, roadmap, challenges, lab, capstone, actions,
// topology, postmortems, chaos, account), each mounting its own routes
// on the mux Routes builds. Config and every helper more than one
// domain package needs live in webcore, the one thing all of them (and
// this package) import, so no domain package needs to import this one
// back.
package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"

	"breaklab/internal/web/webcore"
)

// Config locates everything the local web UI needs to read, and how it
// binds and announces itself. It's an alias for webcore.Config so every
// existing caller (cmd/breaklab/serve.go, tests) keeps working
// unchanged even though Config itself now lives in webcore.
type Config = webcore.Config

// OpenBrowserFunc opens url in the user's default browser.
type OpenBrowserFunc = webcore.OpenBrowserFunc

// Stats and DiskUsage mirror docker.Stats/docker.DiskUsage; see
// webcore's own doc comments. Aliased here so cmd/breaklab (which
// constructs a Config's telemetry funcs) keeps working unchanged.
type Stats = webcore.Stats
type DiskUsage = webcore.DiskUsage

// Serve starts the local web UI and blocks until ctx is cancelled or
// the server fails to serve. Listening happens before anything else, so
// a port already in use produces a specific error naming the port
// instead of a generic bind failure.
func Serve(ctx context.Context, cfg Config) error {
	cfg.Started = time.Now()

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			return fmt.Errorf("port %d is already in use", cfg.Port)
		}
		return err
	}

	if isLoopback(cfg.Host) {
		if cfg.OpenBrowser != nil {
			_ = cfg.OpenBrowser(fmt.Sprintf("http://%s/", addr))
		}
	} else {
		fmt.Fprintf(cfg.Writer(), "warning: the local web UI is unauthenticated and reachable by any device that can reach %s\n", addr)
	}

	srv := &http.Server{Handler: Routes(cfg)}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Serve(ln) }()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

// isLoopback reports whether host is a loopback address (127.0.0.1,
// ::1, or the literal "localhost").
func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
