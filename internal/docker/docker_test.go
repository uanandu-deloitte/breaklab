package docker

import (
	"context"
	"errors"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

type call struct {
	name string
	args []string
}

func fake(calls *[]call, out string, err error) ExecFunc {
	return func(_ context.Context, name string, args []string) ([]byte, error) {
		*calls = append(*calls, call{name, args})
		return []byte(out), err
	}
}

func TestPing(t *testing.T) {
	var calls []call

	if err := NewWithExec(fake(&calls, "", nil)).Ping(context.Background()); err != nil {
		t.Fatalf("success: %v", err)
	}

	err := NewWithExec(fake(&calls, "Cannot connect to the Docker daemon", errors.New("exit status 1"))).Ping(context.Background())
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("daemon down: want ErrUnavailable, got %v", err)
	}

	err = NewWithExec(fake(&calls, "", exec.ErrNotFound)).Ping(context.Background())
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("missing binary: want ErrUnavailable, got %v", err)
	}
}

func TestArgv(t *testing.T) {
	var calls []call
	d := NewWithExec(fake(&calls, "", nil))
	ctx := context.Background()

	_ = d.Build(ctx, "breaklab-web-down", "labs/web-down", "labs/web-down/Dockerfile")
	_ = d.Run(ctx, "breaklab-web-down", "breaklab-web-down", []string{"8080:80"}, false, nil, 0)
	_, _ = d.Exec(ctx, "breaklab-web-down", "bash", "/tmp/v.sh")
	_ = d.Copy(ctx, "labs/web-down/variants/config-typo.sh", "breaklab-web-down", "/tmp/v.sh")
	_ = d.Remove(ctx, "breaklab-web-down")

	want := [][]string{
		{"build", "-t", "breaklab-web-down", "-f", "labs/web-down/Dockerfile", "labs/web-down"},
		{"run", "-d", "--name", "breaklab-web-down", "-p", "8080:80", "breaklab-web-down"},
		{"exec", "breaklab-web-down", "bash", "/tmp/v.sh"},
		{"cp", "labs/web-down/variants/config-typo.sh", "breaklab-web-down:/tmp/v.sh"},
		{"rm", "-f", "breaklab-web-down"},
	}
	if len(calls) != len(want) {
		t.Fatalf("got %d calls, want %d", len(calls), len(want))
	}
	for i, c := range calls {
		if c.name != "docker" || !reflect.DeepEqual(c.args, want[i]) {
			t.Errorf("call %d = %s %v, want docker %v", i, c.name, c.args, want[i])
		}
	}
}

// Privileged must be an explicit, per-call opt-in: false never adds
// --privileged, and true adds it exactly once, in the expected position.
func TestRunPrivilegedFlag(t *testing.T) {
	var calls []call
	d := NewWithExec(fake(&calls, "", nil))
	ctx := context.Background()

	_ = d.Run(ctx, "breaklab-plain", "breaklab-plain", nil, false, nil, 0)
	_ = d.Run(ctx, "breaklab-k8s", "breaklab-k8s", nil, true, nil, 0)

	want := [][]string{
		{"run", "-d", "--name", "breaklab-plain", "breaklab-plain"},
		{"run", "-d", "--name", "breaklab-k8s", "--privileged", "breaklab-k8s"},
	}
	if len(calls) != len(want) {
		t.Fatalf("got %d calls, want %d", len(calls), len(want))
	}
	for i, c := range calls {
		if !reflect.DeepEqual(c.args, want[i]) {
			t.Errorf("call %d args = %v, want %v", i, c.args, want[i])
		}
	}
}

// Volumes and pids-limit are both opt-in: nil/0 never add anything, and
// non-empty/positive values add exactly the expected flags, in order.
func TestRunVolumesAndPidsLimit(t *testing.T) {
	var calls []call
	d := NewWithExec(fake(&calls, "", nil))
	ctx := context.Background()

	_ = d.Run(ctx, "breaklab-plain", "breaklab-plain", nil, false, nil, 0)
	_ = d.Run(ctx, "breaklab-vol", "breaklab-vol", nil, false, []string{"/host/data:/data"}, 0)
	_ = d.Run(ctx, "breaklab-pids", "breaklab-pids", nil, false, nil, 32)
	_ = d.Run(ctx, "breaklab-both", "breaklab-both", nil, false, []string{"/host/a:/a", "/host/b:/b"}, 16)

	want := [][]string{
		{"run", "-d", "--name", "breaklab-plain", "breaklab-plain"},
		{"run", "-d", "--name", "breaklab-vol", "-v", "/host/data:/data", "breaklab-vol"},
		{"run", "-d", "--name", "breaklab-pids", "--pids-limit", "32", "breaklab-pids"},
		{"run", "-d", "--name", "breaklab-both", "-v", "/host/a:/a", "-v", "/host/b:/b", "--pids-limit", "16", "breaklab-both"},
	}
	if len(calls) != len(want) {
		t.Fatalf("got %d calls, want %d", len(calls), len(want))
	}
	for i, c := range calls {
		if !reflect.DeepEqual(c.args, want[i]) {
			t.Errorf("call %d args = %v, want %v", i, c.args, want[i])
		}
	}
}

func TestInvalidNamesRejected(t *testing.T) {
	var calls []call
	d := NewWithExec(fake(&calls, "", nil))
	ctx := context.Background()

	for _, bad := range []string{"a;b", "../x", "Web", "", "a b", "$(id)"} {
		if err := d.Run(ctx, bad, "ok", nil, false, nil, 0); err == nil {
			t.Errorf("Run accepted container name %q", bad)
		}
		if err := d.Build(ctx, bad, ".", "Dockerfile"); err == nil {
			t.Errorf("Build accepted image name %q", bad)
		}
		if err := d.Remove(ctx, bad); err == nil {
			t.Errorf("Remove accepted %q", bad)
		}
		if _, err := d.Exec(ctx, bad, "true"); err == nil {
			t.Errorf("Exec accepted %q", bad)
		}
	}
	if len(calls) != 0 {
		t.Errorf("docker was invoked %d times for invalid names", len(calls))
	}
}

func TestStatus(t *testing.T) {
	var calls []call
	ctx := context.Background()

	t.Run("running container", func(t *testing.T) {
		d := NewWithExec(fake(&calls, "running\tbreaklab-web-down\n", nil))
		state, image, err := d.Status(ctx, "breaklab-web-down")
		if err != nil {
			t.Fatalf("Status: %v", err)
		}
		if state != "running" || image != "breaklab-web-down" {
			t.Errorf("Status = %q, %q, want running, breaklab-web-down", state, image)
		}
	})

	t.Run("stopped container", func(t *testing.T) {
		d := NewWithExec(fake(&calls, "exited\tbreaklab-web-down\n", nil))
		state, _, err := d.Status(ctx, "breaklab-web-down")
		if err != nil {
			t.Fatalf("Status: %v", err)
		}
		if state != "exited" {
			t.Errorf("state = %q, want exited", state)
		}
	})

	t.Run("nonexistent container", func(t *testing.T) {
		d := NewWithExec(fake(&calls, "Error: No such container: breaklab-web-down", errors.New("exit status 1")))
		if _, _, err := d.Status(ctx, "breaklab-web-down"); err == nil {
			t.Fatal("Status on a nonexistent container should return an error")
		}
	})

	t.Run("exec error (e.g. daemon unreachable or ctx cancelled)", func(t *testing.T) {
		d := NewWithExec(fake(&calls, "", context.DeadlineExceeded))
		if _, _, err := d.Status(ctx, "breaklab-web-down"); err == nil {
			t.Fatal("Status should surface an exec error rather than returning a zero value silently")
		}
	})

	t.Run("invalid container name rejected before exec", func(t *testing.T) {
		calls = nil
		d := NewWithExec(fake(&calls, "", nil))
		if _, _, err := d.Status(ctx, "../etc"); err == nil {
			t.Fatal("Status accepted an invalid container name")
		}
		if len(calls) != 0 {
			t.Errorf("docker was invoked for an invalid name")
		}
	})
}

func TestVersion(t *testing.T) {
	var calls []call
	d := NewWithExec(fake(&calls, "26.1.1\n", nil))
	v, err := d.Version(context.Background())
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if v != "26.1.1" {
		t.Errorf("Version = %q, want 26.1.1", v)
	}
	want := []string{"version", "--format", "{{.Server.Version}}"}
	if !reflect.DeepEqual(calls[len(calls)-1].args, want) {
		t.Errorf("args = %v, want %v", calls[len(calls)-1].args, want)
	}
}

func TestStats(t *testing.T) {
	ctx := context.Background()

	t.Run("valid output", func(t *testing.T) {
		var calls []call
		out := `{"CPUPerc":"4.20%","MemUsage":"15.47MiB / 1.952GiB"}` + "\n"
		d := NewWithExec(fake(&calls, out, nil))
		s, err := d.Stats(ctx, "breaklab-web-down")
		if err != nil {
			t.Fatalf("Stats: %v", err)
		}
		if s.CPUPercent != 4.20 {
			t.Errorf("CPUPercent = %v, want 4.20", s.CPUPercent)
		}
		usedMiB, limitGiB := 15.47, 1.952
		wantUsed := uint64(usedMiB * 1024 * 1024)
		if s.MemUsedBytes != wantUsed {
			t.Errorf("MemUsedBytes = %d, want %d", s.MemUsedBytes, wantUsed)
		}
		wantLimit := uint64(limitGiB * 1024 * 1024 * 1024)
		if s.MemLimitBytes != wantLimit {
			t.Errorf("MemLimitBytes = %d, want %d", s.MemLimitBytes, wantLimit)
		}
	})

	t.Run("malformed output", func(t *testing.T) {
		var calls []call
		d := NewWithExec(fake(&calls, "not json", nil))
		if _, err := d.Stats(ctx, "breaklab-web-down"); err == nil {
			t.Fatal("Stats on malformed output should return an error")
		}
	})

	t.Run("invalid container name rejected before exec", func(t *testing.T) {
		var calls []call
		d := NewWithExec(fake(&calls, "", nil))
		if _, err := d.Stats(ctx, "../etc"); err == nil {
			t.Fatal("Stats accepted an invalid container name")
		}
		if len(calls) != 0 {
			t.Errorf("docker was invoked for an invalid name")
		}
	})
}

func TestLogs(t *testing.T) {
	var calls []call
	d := NewWithExec(fake(&calls, "line one\nline two\n", nil))
	out, err := d.Logs(context.Background(), "breaklab-web-down", 20)
	if err != nil {
		t.Fatalf("Logs: %v", err)
	}
	if string(out) != "line one\nline two\n" {
		t.Errorf("Logs output = %q", out)
	}
	want := []string{"logs", "--tail", "20", "breaklab-web-down"}
	if !reflect.DeepEqual(calls[len(calls)-1].args, want) {
		t.Errorf("args = %v, want %v", calls[len(calls)-1].args, want)
	}
}

func TestLogsFollowClosesProcess(t *testing.T) {
	// A real but short-lived process: sleep in the background so we can
	// observe that Close() kills it rather than letting it run to
	// completion on its own.
	d := New()
	rc, err := realStream(context.Background(), "sleep", []string{"5"})
	if err != nil {
		t.Fatalf("realStream: %v", err)
	}
	ps, ok := rc.(*procStream)
	if !ok {
		t.Fatalf("realStream returned %T, want *procStream", rc)
	}
	if err := rc.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	state, err := ps.cmd.Process.Wait()
	if err != nil {
		// Already reaped by the internal cmd.Wait() goroutine — that's fine,
		// it just means Close()'s kill raced a wait that already happened.
		return
	}
	if state.ExitCode() == 0 {
		t.Errorf("process exited cleanly (code 0); want it killed")
	}
	_ = d // keep d referenced for symmetry with other tests in this file
}

func TestDiskUsage(t *testing.T) {
	ctx := context.Background()

	t.Run("valid output", func(t *testing.T) {
		var calls []call
		out := `{"Type":"Images","TotalCount":"3","Active":"1","Size":"1.2GB","Reclaimable":"0B (0%)"}` + "\n" +
			`{"Type":"Containers","TotalCount":"1","Active":"1","Size":"0B","Reclaimable":"0B"}` + "\n"
		d := NewWithExec(fake(&calls, out, nil))
		du, err := d.DiskUsage(ctx)
		if err != nil {
			t.Fatalf("DiskUsage: %v", err)
		}
		sizeGB := 1.2
		want := uint64(sizeGB * 1000 * 1000 * 1000)
		if du.ImagesSizeBytes != want {
			t.Errorf("ImagesSizeBytes = %d, want %d", du.ImagesSizeBytes, want)
		}
	})

	t.Run("no Images row", func(t *testing.T) {
		var calls []call
		d := NewWithExec(fake(&calls, `{"Type":"Containers","Size":"0B"}`+"\n", nil))
		if _, err := d.DiskUsage(ctx); err == nil {
			t.Fatal("DiskUsage with no Images row should return an error")
		}
	})
}

func TestRemoveIdempotent(t *testing.T) {
	var calls []call
	d := NewWithExec(fake(&calls, "Error response from daemon: No such container: x", errors.New("exit status 1")))
	if err := d.Remove(context.Background(), "x"); err != nil {
		t.Fatalf("removing a missing container should succeed: %v", err)
	}

	d = NewWithExec(fake(&calls, "permission denied", errors.New("exit status 1")))
	if err := d.Remove(context.Background(), "x"); err == nil || !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("other errors should surface, got %v", err)
	}
}
