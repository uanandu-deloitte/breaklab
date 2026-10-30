// Package webcore holds the pieces internal/web's domain packages all
// need to share — Config and the handful of helpers that would
// otherwise force every domain package to import the composition root
// (and the root to import them back), which Go doesn't allow.
package webcore

import (
	"io"
	"os"
	"time"

	"breaklab/internal/account"
	"breaklab/internal/runner"
)

// OpenBrowserFunc opens url in the user's default browser.
type OpenBrowserFunc func(url string) error

// Config locates everything the local web UI needs to read, and how it
// binds and announces itself.
type Config struct {
	LabsDir      string
	RoadmapPath  string
	SessionPath  string
	HistoryPath  string
	ProgressPath string

	// TemplatesDir now only locates vendored static assets (fonts,
	// cytoscape.min.js) — every page template that used to live here
	// migrated to the `frontend` capability (FrontendDir below).
	TemplatesDir string

	// FrontendDir is the Astro build's static output (frontend/dist/),
	// serving every page (the `frontend` capability). Empty mounts no
	// frontend routes at all — only used by tests that don't need any
	// page route; `breaklab serve` always sets it (see --frontend's
	// default).
	FrontendDir string

	// PostmortemsDir confines which files a progress.json postmortem path
	// may name; anything outside it is never linked or read.
	PostmortemsDir string

	// Docker runs the lab lifecycle (start/check/hint/stop) for the lab
	// detail page's action routes. Nil (e.g. most test Configs, which
	// only exercise the read-only views) makes every action route fail
	// with a clear error instead of a nil-pointer panic.
	Docker runner.Docker

	Host string
	Port int

	// OpenBrowser opens the server's URL once it's listening, but only
	// when Host is a loopback address (see Serve). Nil means don't open
	// one — used by tests and by any caller that handles it itself.
	OpenBrowser OpenBrowserFunc

	// DockerStatus looks up a container's live status for the Topology
	// view. Nil (e.g. in most test Configs) degrades every container to
	// status "unknown" rather than failing.
	DockerStatus StatusFunc

	// DockerVersion, ContainerStats, ContainerLogs, ContainerLogsFollow,
	// and DiskUsage back the dashboard/lab-detail/challenge-browser's
	// live telemetry panels, following the same nil-degrades-gracefully
	// contract as DockerStatus — never required, never a 500.
	DockerVersion       VersionFunc
	ContainerStats      StatsFunc
	ContainerLogs       LogsFunc
	ContainerLogsFollow LogsFollowFunc
	DiskUsage           DiskUsageFunc

	// Started is when this server began listening, for the dashboard's
	// uptime display. Left zero in most test Configs, which just shows
	// no uptime.
	Started time.Time

	// Out receives the one-time non-loopback warning; defaults to
	// os.Stderr when nil.
	Out io.Writer

	// AccountStore backs the optional hosted-account layer (signup,
	// login, sync, public profile). Nil — the default, and every
	// Config that existed before this field — means none of it runs:
	// every account-aware route or branch is mounted or taken only
	// when this is non-nil, so local-only usage is unaffected.
	AccountStore account.Store
}

// Writer returns where the server should print warnings — cfg.Out, or
// os.Stderr when unset.
func (cfg Config) Writer() io.Writer {
	if cfg.Out != nil {
		return cfg.Out
	}
	return os.Stderr
}
