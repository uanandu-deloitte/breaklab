package webcore

import (
	"context"
	"io"
	"time"
)

// DockerStatusTimeout bounds how long a Docker-backed view waits before
// treating it as unreachable — every other local UI route works
// entirely from files on disk and never blocks on an external process.
const DockerStatusTimeout = 2 * time.Second

// StatusFunc looks up a container's live running state (Docker's own
// value, e.g. "running" or "exited") and image reference. *docker.Docker's
// Status method satisfies this; Config.DockerStatus holds it so this
// package never needs to import internal/docker, matching the
// OpenBrowserFunc pattern. A nil StatusFunc (e.g. in a test Config that
// doesn't set one) is treated the same as a lookup error: unknown.
type StatusFunc func(ctx context.Context, name string) (state, image string, err error)

// Stats mirrors docker.Stats without importing internal/docker, the same
// reasoning StatusFunc already documents: this package never needs to
// know Docker is the thing behind these funcs.
type Stats struct {
	CPUPercent    float64
	MemUsedBytes  uint64
	MemLimitBytes uint64
}

// DiskUsage mirrors docker.DiskUsage, same reasoning as Stats.
type DiskUsage struct {
	ImagesSizeBytes uint64
}

// VersionFunc returns the Docker server's version string for the
// dashboard header. Nil degrades to an empty string (shown as "—").
type VersionFunc func(ctx context.Context) (string, error)

// StatsFunc is a read-only, single-sample resource lookup for one
// container, used by the dashboard and lab detail's sandbox panels. Nil
// degrades to a zero Stats (shown as "—").
type StatsFunc func(ctx context.Context, container string) (Stats, error)

// LogsFunc returns the last tailLines lines of a container's own
// combined stdout/stderr, for the dashboard/lab-detail diagnostic
// output panel. Nil degrades to no panel being shown.
type LogsFunc func(ctx context.Context, container string, tailLines int) ([]byte, error)

// LogsFollowFunc streams a container's own combined stdout/stderr as
// it's produced, for the lab detail page's live log tail. Nil means the
// live-tail route degrades to not-found rather than panicking.
type LogsFollowFunc func(ctx context.Context, container string) (io.ReadCloser, error)

// DiskUsageFunc reports Docker's image-store disk usage, for the
// challenge browser's sandbox-image-cache tile. Nil degrades to a zero
// DiskUsage (shown as "—").
type DiskUsageFunc func(ctx context.Context) (DiskUsage, error)
