// Package docker is a thin wrapper around the docker CLI. Every call is an
// os/exec invocation with an argument slice; nothing goes through a shell.
package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ErrUnavailable means the docker CLI is missing or the daemon is unreachable.
var ErrUnavailable = errors.New("docker is not available")

var namePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

// ExecFunc runs a command to completion and returns its combined stdout
// and stderr. Every docker subcommand that terminates on its own (build,
// run, stats --no-stream, ...) goes through this.
type ExecFunc func(ctx context.Context, name string, args []string) ([]byte, error)

// StreamFunc starts a command and returns a ReadCloser over its combined
// output as the command produces it. Close must terminate the process if
// it hasn't exited on its own. Only LogsFollow uses this — every other
// docker call runs to completion and fits ExecFunc instead.
type StreamFunc func(ctx context.Context, name string, args []string) (io.ReadCloser, error)

// Docker drives the docker CLI through injectable command runners.
type Docker struct {
	exec   ExecFunc
	stream StreamFunc
}

// New returns a Docker that runs the real docker binary.
func New() *Docker {
	return &Docker{exec: realExec, stream: realStream}
}

// NewWithExec returns a Docker that runs non-streaming commands through
// fn; used by tests. Streaming (LogsFollow) still runs for real — use
// NewWithExecAndStream to fake that too.
func NewWithExec(fn ExecFunc) *Docker {
	return &Docker{exec: fn, stream: realStream}
}

// NewWithExecAndStream returns a Docker that runs commands through fn and
// sfn; used by tests that exercise LogsFollow.
func NewWithExecAndStream(fn ExecFunc, sfn StreamFunc) *Docker {
	return &Docker{exec: fn, stream: sfn}
}

func realExec(ctx context.Context, name string, args []string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// realStream starts the command with its combined stdout+stderr piped
// through the returned ReadCloser. Close kills the process if it's still
// running (the caller's only way to stop a "docker logs -f" that follows
// forever) and always releases the pipe.
func realStream(ctx context.Context, name string, args []string) (io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	if err := cmd.Start(); err != nil {
		pw.Close()
		return nil, err
	}
	go func() { pw.CloseWithError(cmd.Wait()) }()
	return &procStream{r: pr, cmd: cmd}, nil
}

type procStream struct {
	r   *io.PipeReader
	cmd *exec.Cmd
}

func (p *procStream) Read(b []byte) (int, error) { return p.r.Read(b) }

func (p *procStream) Close() error {
	if p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}
	return p.r.Close()
}

func (d *Docker) docker(ctx context.Context, args ...string) ([]byte, error) {
	return d.exec(ctx, "docker", args)
}

func (d *Docker) dockerStream(ctx context.Context, args ...string) (io.ReadCloser, error) {
	return d.stream(ctx, "docker", args)
}

func validName(kind, name string) error {
	if !namePattern.MatchString(name) {
		return fmt.Errorf("invalid %s %q: must match ^[a-z0-9-]+$", kind, name)
	}
	return nil
}

// Ping reports ErrUnavailable when the CLI is missing or the daemon is down.
func (d *Docker) Ping(ctx context.Context) error {
	out, err := d.docker(ctx, "info")
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, strings.TrimSpace(string(out)+" "+err.Error()))
	}
	return nil
}

// Build builds image from the Dockerfile at dockerfile inside contextDir.
func (d *Docker) Build(ctx context.Context, image, contextDir, dockerfile string) error {
	if err := validName("image name", image); err != nil {
		return err
	}
	out, err := d.docker(ctx, "build", "-t", image, "-f", dockerfile, contextDir)
	if err != nil {
		return fmt.Errorf("docker build failed: %w\n%s", err, tail(out))
	}
	return nil
}

// Run starts a detached container publishing each "host:container" port.
// privileged, when true, adds --privileged — only a lab that explicitly
// declares environment.privileged: true in its lab.yaml reaches this as
// true; every other caller passes false. volumes is a list of
// already-resolved "hostAbsPath:containerPath" bind mounts (the caller
// owns creating and cleaning up the host side); pidsLimit, when > 0,
// adds --pids-limit.
func (d *Docker) Run(ctx context.Context, name, image string, ports []string, privileged bool, volumes []string, pidsLimit int) error {
	if err := validName("container name", name); err != nil {
		return err
	}
	if err := validName("image name", image); err != nil {
		return err
	}
	args := []string{"run", "-d", "--name", name}
	if privileged {
		args = append(args, "--privileged")
	}
	for _, p := range ports {
		args = append(args, "-p", p)
	}
	for _, v := range volumes {
		args = append(args, "-v", v)
	}
	if pidsLimit > 0 {
		args = append(args, "--pids-limit", strconv.Itoa(pidsLimit))
	}
	args = append(args, image)
	out, err := d.docker(ctx, args...)
	if err != nil {
		return fmt.Errorf("docker run failed: %w\n%s", err, tail(out))
	}
	return nil
}

// Exec runs argv inside a running container as its default user (root).
func (d *Docker) Exec(ctx context.Context, container string, argv ...string) ([]byte, error) {
	if err := validName("container name", container); err != nil {
		return nil, err
	}
	args := append([]string{"exec", container}, argv...)
	return d.docker(ctx, args...)
}

// Copy copies the host file src to dst inside the container.
func (d *Docker) Copy(ctx context.Context, src, container, dst string) error {
	if err := validName("container name", container); err != nil {
		return err
	}
	out, err := d.docker(ctx, "cp", src, container+":"+dst)
	if err != nil {
		return fmt.Errorf("docker cp failed: %w\n%s", err, tail(out))
	}
	return nil
}

// Remove force-removes a container; a container that doesn't exist is fine.
func (d *Docker) Remove(ctx context.Context, name string) error {
	if err := validName("container name", name); err != nil {
		return err
	}
	out, err := d.docker(ctx, "rm", "-f", name)
	if err != nil && !strings.Contains(string(out), "No such container") {
		return fmt.Errorf("docker rm failed: %w\n%s", err, tail(out))
	}
	return nil
}

// VolumeRemove force-removes a named volume; a volume that doesn't exist
// is fine. Used to clear a previous attempt's volume before recreating
// it, the same "wipe before reuse" semantics as a bind-mounted host dir.
func (d *Docker) VolumeRemove(ctx context.Context, name string) error {
	if err := validName("volume name", name); err != nil {
		return err
	}
	out, err := d.docker(ctx, "volume", "rm", "-f", name)
	if err != nil && !strings.Contains(string(out), "No such volume") {
		return fmt.Errorf("docker volume rm failed: %w\n%s", err, tail(out))
	}
	return nil
}

// VolumeCreate creates a named (daemon-managed) volume. Unlike a host
// bind mount, this lives entirely inside the Docker VM's own filesystem,
// so it gets real multi-UID permission enforcement even on backends
// (e.g. Docker Desktop's default Mac file sharing) where a host bind
// mount doesn't honor chown to a non-root uid.
func (d *Docker) VolumeCreate(ctx context.Context, name string) error {
	if err := validName("volume name", name); err != nil {
		return err
	}
	out, err := d.docker(ctx, "volume", "create", name)
	if err != nil {
		return fmt.Errorf("docker volume create failed: %w\n%s", err, tail(out))
	}
	return nil
}

// Exists reports whether a container with this name exists.
func (d *Docker) Exists(ctx context.Context, name string) (bool, error) {
	if err := validName("container name", name); err != nil {
		return false, err
	}
	_, err := d.docker(ctx, "container", "inspect", name)
	return err == nil, nil
}

// Status is a read-only lookup of a container's running state (Docker's
// own value, e.g. "running" or "exited") and image reference. An error
// means the container doesn't exist, the daemon is unreachable, or ctx
// was cancelled (e.g. a caller-imposed timeout) — the caller decides how
// to degrade; this method never blocks longer than ctx allows.
func (d *Docker) Status(ctx context.Context, name string) (state, image string, err error) {
	if err := validName("container name", name); err != nil {
		return "", "", err
	}
	out, err := d.docker(ctx, "container", "inspect", "--format", "{{.State.Status}}\t{{.Config.Image}}", name)
	if err != nil {
		return "", "", fmt.Errorf("docker container inspect failed: %w\n%s", err, tail(out))
	}
	fields := strings.SplitN(strings.TrimSpace(string(out)), "\t", 2)
	if len(fields) != 2 {
		return "", "", fmt.Errorf("unexpected docker inspect output: %q", tail(out))
	}
	return fields[0], fields[1], nil
}

func tail(out []byte) string {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 15 {
		lines = lines[len(lines)-15:]
	}
	return strings.Join(lines, "\n")
}

// Version returns the Docker server's version string (e.g. "26.1.1").
func (d *Docker) Version(ctx context.Context) (string, error) {
	out, err := d.docker(ctx, "version", "--format", "{{.Server.Version}}")
	if err != nil {
		return "", fmt.Errorf("docker version failed: %w\n%s", err, tail(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// Stats is a point-in-time resource snapshot for one running container.
type Stats struct {
	CPUPercent    float64
	MemUsedBytes  uint64
	MemLimitBytes uint64
}

type statsJSON struct {
	CPUPerc  string `json:"CPUPerc"`
	MemUsage string `json:"MemUsage"`
}

// Stats is a read-only, single-sample resource snapshot (`docker stats
// --no-stream`) for container. An error means the container doesn't
// exist, the daemon is unreachable, output didn't parse, or ctx was
// cancelled — the caller decides how to degrade.
func (d *Docker) Stats(ctx context.Context, container string) (Stats, error) {
	if err := validName("container name", container); err != nil {
		return Stats{}, err
	}
	out, err := d.docker(ctx, "stats", "--no-stream", "--format", "{{json .}}", container)
	if err != nil {
		return Stats{}, fmt.Errorf("docker stats failed: %w\n%s", err, tail(out))
	}

	var raw statsJSON
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(out))), &raw); err != nil {
		return Stats{}, fmt.Errorf("unexpected docker stats output: %q: %w", tail(out), err)
	}
	cpu, err := parsePercent(raw.CPUPerc)
	if err != nil {
		return Stats{}, fmt.Errorf("unexpected docker stats CPUPerc %q: %w", raw.CPUPerc, err)
	}
	used, limit, err := parseMemUsage(raw.MemUsage)
	if err != nil {
		return Stats{}, fmt.Errorf("unexpected docker stats MemUsage %q: %w", raw.MemUsage, err)
	}
	return Stats{CPUPercent: cpu, MemUsedBytes: used, MemLimitBytes: limit}, nil
}

// Logs returns the last tailLines lines of container's own combined
// stdout/stderr.
func (d *Docker) Logs(ctx context.Context, container string, tailLines int) ([]byte, error) {
	if err := validName("container name", container); err != nil {
		return nil, err
	}
	out, err := d.docker(ctx, "logs", "--tail", strconv.Itoa(tailLines), container)
	if err != nil {
		return nil, fmt.Errorf("docker logs failed: %w\n%s", err, tail(out))
	}
	return out, nil
}

// LogsFollow streams container's own combined stdout/stderr as it's
// produced, starting from its last 20 lines. The caller must Close the
// returned ReadCloser (e.g. when the request that wanted it ends) —
// Close terminates the underlying "docker logs -f" process rather than
// leaving it running.
func (d *Docker) LogsFollow(ctx context.Context, container string) (io.ReadCloser, error) {
	if err := validName("container name", container); err != nil {
		return nil, err
	}
	return d.dockerStream(ctx, "logs", "-f", "--tail", "20", container)
}

// DiskUsage is the disk space Docker's image store is using.
type DiskUsage struct {
	ImagesSizeBytes uint64
}

type diskUsageRowJSON struct {
	Type string `json:"Type"`
	Size string `json:"Size"`
}

// DiskUsage reports how much disk space Docker's images are using
// (`docker system df`), for the challenge browser's sandbox-image-cache
// tile.
func (d *Docker) DiskUsage(ctx context.Context) (DiskUsage, error) {
	out, err := d.docker(ctx, "system", "df", "--format", "{{json .}}")
	if err != nil {
		return DiskUsage{}, fmt.Errorf("docker system df failed: %w\n%s", err, tail(out))
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var row diskUsageRowJSON
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return DiskUsage{}, fmt.Errorf("unexpected docker system df output: %q: %w", tail(out), err)
		}
		if row.Type != "Images" {
			continue
		}
		size, err := parseHumanSize(row.Size)
		if err != nil {
			return DiskUsage{}, fmt.Errorf("unexpected docker system df Images size %q: %w", row.Size, err)
		}
		return DiskUsage{ImagesSizeBytes: size}, nil
	}
	return DiskUsage{}, fmt.Errorf("docker system df output had no Images row: %q", tail(out))
}

func parsePercent(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(s), "%"), 64)
}

func parseMemUsage(s string) (used, limit uint64, err error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected \"used / limit\" format")
	}
	used, err = parseHumanSize(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	limit, err = parseHumanSize(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	return used, limit, nil
}

// sizePattern matches docker's human-readable byte sizes: decimal units
// from `docker system df` (e.g. "1.2GB") and binary units from `docker
// stats` (e.g. "15.47MiB").
var sizePattern = regexp.MustCompile(`(?i)^([0-9]*\.?[0-9]+)\s*([kmgt]?)(i?)b$`)

func parseHumanSize(s string) (uint64, error) {
	m := sizePattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0, fmt.Errorf("unrecognized size %q", s)
	}
	value, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, err
	}
	base := 1000.0
	if strings.EqualFold(m[3], "i") {
		base = 1024.0
	}
	switch strings.ToLower(m[2]) {
	case "k":
		value *= base
	case "m":
		value *= base * base
	case "g":
		value *= base * base * base
	case "t":
		value *= base * base * base * base
	}
	return uint64(value), nil
}
