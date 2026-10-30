package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"breaklab/internal/account"
	"breaklab/internal/docker"
	"breaklab/internal/web"
)

// serveFlags defines `breaklab serve`'s flags. The host defaults to
// loopback, so binding anywhere reachable from the LAN always takes an
// explicit --host. --db is a separate, opt-in deployment shape: with
// it, the server also mounts hosted-account routes (signup, login,
// sync, public profile); without it (the default), none of that code
// runs at all, same as before this flag existed. --frontend defaults
// to frontend/dist (not opt-in like --db): once a page has migrated to
// the `frontend` capability, there is no html/template fallback for
// it, so serving that page at all requires the Astro build to exist —
// the same one-time prerequisite `go build` already is.
func serveFlags() (fs *flag.FlagSet, host *string, port *int, db *string, frontendDir *string) {
	fs = flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: breaklab serve [--host 127.0.0.1] [--port 8787] [--db path] [--frontend path]")
		fmt.Fprintln(os.Stderr, "\nServes the local dashboard, roadmap, challenge browser, lab detail, and postmortem viewer.")
		fmt.Fprintln(os.Stderr, "--db opts into hosted accounts (signup/login/sync/public profile); omit it for the default, account-free local server.")
		fmt.Fprintln(os.Stderr, "--frontend opts into the Astro-built pages (see frontend/); omit it to keep every page on its existing html/template rendering.")
		fs.PrintDefaults()
	}
	host = fs.String("host", "127.0.0.1", "address to bind (loopback by default; a non-loopback host has no authentication)")
	port = fs.Int("port", 8787, "port to bind")
	db = fs.String("db", "", "path to a SQLite database file; opts into hosted accounts when set (unset by default)")
	frontendDir = fs.String("frontend", filepath.Join("frontend", "dist"), "path to the Astro build's static output; migrated pages (e.g. the dashboard) need it built once via `npm run build` in frontend/")
	return fs, host, port, db, frontendDir
}

func runServe(args []string) int {
	fs, host, port, db, frontendDir := serveFlags()
	if done, code := parseFlags(fs, args); done {
		return code
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return 2
	}

	ctx, cancel := commandContext()
	defer cancel()

	fmt.Printf("Serving BreakLab at http://%s:%d/\n", *host, *port)

	d := docker.New()
	cfg := web.Config{
		LabsDir:        "labs",
		RoadmapPath:    "roadmap.yaml",
		SessionPath:    filepath.Join(".breaklab", "session.json"),
		HistoryPath:    filepath.Join(".breaklab", "history.json"),
		ProgressPath:   filepath.Join("progress", "progress.json"),
		TemplatesDir:   filepath.Join("templates", "web"),
		FrontendDir:    *frontendDir,
		PostmortemsDir: "postmortems",
		Host:           *host,
		Port:           *port,
		OpenBrowser:    web.DefaultOpenBrowser,
		DockerStatus:   d.Status,
		Docker:         d,

		DockerVersion: d.Version,
		ContainerStats: func(ctx context.Context, container string) (web.Stats, error) {
			s, err := d.Stats(ctx, container)
			return web.Stats(s), err
		},
		ContainerLogs:       d.Logs,
		ContainerLogsFollow: d.LogsFollow,
		DiskUsage: func(ctx context.Context) (web.DiskUsage, error) {
			du, err := d.DiskUsage(ctx)
			return web.DiskUsage(du), err
		},
	}

	if *db != "" {
		store, err := account.Open(*db)
		if err != nil {
			return report(err)
		}
		defer store.Close()
		cfg.AccountStore = store
	}

	if err := web.Serve(ctx, cfg); err != nil {
		return report(err)
	}
	return 0
}
