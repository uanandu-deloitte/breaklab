package runner

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"breaklab/internal/docker"
	"breaklab/internal/progress"
)

type fakeDocker struct {
	calls   []string
	pingErr error
	failOn  string
}

func (f *fakeDocker) rec(s string) error {
	f.calls = append(f.calls, s)
	if f.failOn != "" && strings.HasPrefix(s, f.failOn) {
		return errors.New("boom")
	}
	return nil
}
func (f *fakeDocker) Ping(context.Context) error { return f.pingErr }
func (f *fakeDocker) Build(_ context.Context, image, _, _ string) error {
	return f.rec("build " + image)
}
func (f *fakeDocker) Run(_ context.Context, name, _ string, ports []string, privileged bool, volumes []string, pidsLimit int) error {
	suffix := ""
	if privileged {
		suffix = " privileged"
	}
	if len(volumes) > 0 {
		suffix += " vol=" + strings.Join(volumes, ",")
	}
	if pidsLimit > 0 {
		suffix += " pids=" + strconv.Itoa(pidsLimit)
	}
	return f.rec("run " + name + " " + strings.Join(ports, ",") + suffix)
}
func (f *fakeDocker) Exec(_ context.Context, c string, argv ...string) ([]byte, error) {
	return nil, f.rec("exec " + c + " " + strings.Join(argv, " "))
}
func (f *fakeDocker) Copy(_ context.Context, _, c, dst string) error {
	return f.rec("cp " + c + " " + dst)
}
func (f *fakeDocker) Remove(_ context.Context, n string) error { return f.rec("rm " + n) }
func (f *fakeDocker) VolumeRemove(_ context.Context, n string) error {
	return f.rec("volume-rm " + n)
}
func (f *fakeDocker) VolumeCreate(_ context.Context, n string) error {
	return f.rec("volume-create " + n)
}

const roadmapYAML = "stages:\n  - id: linux\n    title: Linux fundamentals\n"

func labYAML(id, runIn string) string {
	return "id: " + id + `
title: Test lab
stage: linux
difficulty: 1
estimated_minutes: 5
environment:
  dockerfile: Dockerfile
  ports: ["8080:80"]
check:
  script: check.sh
  run_in: ` + runIn + `
hints:
  - general hint one
  - general hint two
variants:
  - id: first-variant
    script: variants/first.sh
    weight: 1
    hints:
      - first-variant hint
  - id: second-variant
    script: variants/second.sh
    weight: 1
    hints:
      - second-variant hint
`
}

func writeLab(t *testing.T, root, id, runIn string) {
	t.Helper()
	dir := filepath.Join(root, "labs", id)
	if err := os.MkdirAll(filepath.Join(dir, "variants"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"lab.yaml":           labYAML(id, runIn),
		"Dockerfile":         "FROM scratch\n",
		"TICKET.md":          "# INC: " + id + " is down\n",
		"check.sh":           "#!/usr/bin/env bash\necho PASS: fine\n",
		"variants/first.sh":  "#!/usr/bin/env bash\n",
		"variants/second.sh": "#!/usr/bin/env bash\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func setup(t *testing.T, labs ...string) (*Runner, *fakeDocker, *bytes.Buffer) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "roadmap.yaml"), []byte(roadmapYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, id := range labs {
		writeLab(t, root, id, "host")
	}
	fd := &fakeDocker{}
	out := &bytes.Buffer{}
	r := &Runner{
		Docker:         fd,
		LabsDir:        filepath.Join(root, "labs"),
		RoadmapPath:    filepath.Join(root, "roadmap.yaml"),
		SessionPath:    filepath.Join(root, ".breaklab", "session.json"),
		HistoryPath:    filepath.Join(root, ".breaklab", "history.json"),
		ProgressPath:   filepath.Join(root, "progress", "progress.json"),
		TemplatePath:   filepath.Join("..", "..", "templates", "postmortem.md"),
		PostmortemsDir: filepath.Join(root, "postmortems"),
		In:             strings.NewReader(""),
		Out:            out,
		HostExec: func(context.Context, string, string, []string, []string) ([]byte, error) {
			return []byte("PASS: fine\n"), nil
		},
		OpenEditor: func(context.Context, string) error { return nil },
	}
	// Most tests here care about start/check/reset/stop plumbing, not
	// variant selection, so pin it to a known value; tests about selection
	// itself (BREAKLAB_FORCE_VARIANT, history) unset this explicitly.
	t.Setenv(forceVariantEnv, "first-variant")
	return r, fd, out
}

func TestStartCallOrderAndHiddenVariant(t *testing.T) {
	r, fd, out := setup(t, "alpha")

	if err := r.Start(context.Background(), "alpha", false); err != nil {
		t.Fatal(err)
	}

	want := []string{
		"rm breaklab-alpha",
		"build breaklab-alpha",
		"run breaklab-alpha 8080:80",
		"cp breaklab-alpha " + containerVariantPath,
		"exec breaklab-alpha bash " + containerVariantPath,
		"exec breaklab-alpha rm -f " + containerVariantPath,
	}
	if !reflect.DeepEqual(fd.calls, want) {
		t.Errorf("calls = %v\nwant    %v", fd.calls, want)
	}

	if strings.Contains(out.String(), "first-variant") {
		t.Errorf("start printed the variant id: %q", out.String())
	}
	if !strings.Contains(out.String(), "INC: alpha is down") || !strings.Contains(out.String(), "docker exec -it breaklab-alpha bash") {
		t.Errorf("start output missing ticket or shell command: %q", out.String())
	}

	s, err := LoadSession(r.SessionPath)
	if err != nil || s == nil || s.Lab != "alpha" || s.Variant != "first-variant" {
		t.Errorf("session = %+v, err %v", s, err)
	}
}

func TestStartCleansUpOnFailure(t *testing.T) {
	r, fd, _ := setup(t, "alpha")
	fd.failOn = "exec breaklab-alpha bash"

	if err := r.Start(context.Background(), "alpha", false); err == nil {
		t.Fatal("expected error")
	}
	if last := fd.calls[len(fd.calls)-1]; last != "rm breaklab-alpha" {
		t.Errorf("last call = %q, want container removal", last)
	}
	if s, _ := LoadSession(r.SessionPath); s != nil {
		t.Errorf("session written despite failure: %+v", s)
	}
}

// TestStartResolvesVolumes covers resolveVolumes end to end: a declared
// volume becomes a real, freshly-emptied host directory, the Run call
// carries its "hostAbsPath:containerPath" mount and the lab's
// pids_limit, and a file left over from a previous attempt doesn't
// survive the next Start.
func TestStartResolvesVolumes(t *testing.T) {
	r, fd, _ := setup(t, "alpha")

	labYAMLPath := filepath.Join(r.LabsDir, "alpha", "lab.yaml")
	data, err := os.ReadFile(labYAMLPath)
	if err != nil {
		t.Fatal(err)
	}
	withVolume := strings.Replace(string(data), "environment:\n  dockerfile: Dockerfile\n  ports: [\"8080:80\"]\n",
		"environment:\n  dockerfile: Dockerfile\n  ports: [\"8080:80\"]\n  pids_limit: 32\n  volumes:\n    - name: data\n      container_path: /data\n", 1)
	if withVolume == string(data) {
		t.Fatal("test fixture's environment block didn't match; update the replacement string")
	}
	if err := os.WriteFile(labYAMLPath, []byte(withVolume), 0o644); err != nil {
		t.Fatal(err)
	}

	volDir := filepath.Join(filepath.Dir(r.SessionPath), "volumes", "alpha", "data")
	if err := os.MkdirAll(volDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(volDir, "leftover-from-last-attempt.txt")
	if err := os.WriteFile(stale, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := r.Start(context.Background(), "alpha", false); err != nil {
		t.Fatal(err)
	}

	absVolDir, err := filepath.Abs(volDir)
	if err != nil {
		t.Fatal(err)
	}
	wantRun := "run breaklab-alpha 8080:80 vol=" + absVolDir + ":/data pids=32"
	var gotRun string
	for _, c := range fd.calls {
		if strings.HasPrefix(c, "run ") {
			gotRun = c
		}
	}
	if gotRun != wantRun {
		t.Errorf("run call = %q, want %q", gotRun, wantRun)
	}

	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Errorf("leftover file from a previous attempt survived Start: %v", err)
	}
	if info, err := os.Stat(volDir); err != nil || !info.IsDir() {
		t.Errorf("expected volume directory to still exist, freshly emptied: %v", err)
	}
}

// TestStartRefusesLockedCapstone covers checkCapstoneUnlocked directly:
// a lab under a stage marked capstone: true refuses to Start until every
// regular stage's labs are solved, and the gate is checked before any
// Docker call (no half-started container on a locked attempt).
func TestStartRefusesLockedCapstone(t *testing.T) {
	root := t.TempDir()
	roadmap := "stages:\n  - id: linux\n    title: Linux fundamentals\n  - id: capstone\n    title: Capstone\n    capstone: true\n"
	if err := os.WriteFile(filepath.Join(root, "roadmap.yaml"), []byte(roadmap), 0o644); err != nil {
		t.Fatal(err)
	}
	writeLab(t, root, "alpha", "host")
	writeLab(t, root, "finale", "host")
	finaleYAML := strings.Replace(labYAML("finale", "host"), "stage: linux", "stage: capstone", 1)
	if finaleYAML == labYAML("finale", "host") {
		t.Fatal("test fixture's stage line didn't match; update the replacement string")
	}
	if err := os.WriteFile(filepath.Join(root, "labs", "finale", "lab.yaml"), []byte(finaleYAML), 0o644); err != nil {
		t.Fatal(err)
	}

	fd := &fakeDocker{}
	r := &Runner{
		Docker:       fd,
		LabsDir:      filepath.Join(root, "labs"),
		RoadmapPath:  filepath.Join(root, "roadmap.yaml"),
		SessionPath:  filepath.Join(root, ".breaklab", "session.json"),
		HistoryPath:  filepath.Join(root, ".breaklab", "history.json"),
		ProgressPath: filepath.Join(root, "progress", "progress.json"),
		Out:          &bytes.Buffer{},
	}
	t.Setenv(forceVariantEnv, "first-variant")

	if err := r.Start(context.Background(), "finale", false); err == nil {
		t.Fatal("expected locked capstone to refuse Start")
	}
	if len(fd.calls) != 0 {
		t.Errorf("docker used before the capstone gate was checked: %v", fd.calls)
	}

	index, err := progress.AppendActive(r.ProgressPath, "alpha", "first-variant")
	if err != nil {
		t.Fatal(err)
	}
	if err := progress.MarkSolved(r.ProgressPath, index, 0); err != nil {
		t.Fatal(err)
	}

	if err := r.Start(context.Background(), "finale", false); err != nil {
		t.Fatalf("expected unlocked capstone to start: %v", err)
	}
}

func TestStartDockerUnavailable(t *testing.T) {
	r, fd, _ := setup(t, "alpha")
	fd.pingErr = docker.ErrUnavailable

	err := r.Start(context.Background(), "alpha", false)
	if !errors.Is(err, docker.ErrUnavailable) {
		t.Fatalf("err = %v", err)
	}
	if len(fd.calls) != 0 {
		t.Errorf("docker used after failed ping: %v", fd.calls)
	}
}

func TestStartUnknownLab(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	if err := r.Start(context.Background(), "nope", false); err == nil || !strings.Contains(err.Error(), "unknown lab") {
		t.Fatalf("err = %v", err)
	}
}

func TestSwitchLabConfirmation(t *testing.T) {
	cases := []struct {
		name        string
		interactive bool
		input       string
		yes         bool
		wantSwitch  bool
	}{
		{"answer yes", true, "y\n", false, true},
		{"answer no", true, "n\n", false, false},
		{"non-interactive", false, "y\n", false, false},
		{"--yes flag", false, "", true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, fd, _ := setup(t, "alpha", "beta")
			if err := r.Start(context.Background(), "alpha", false); err != nil {
				t.Fatal(err)
			}
			fd.calls = nil
			r.Interactive = tc.interactive
			r.In = strings.NewReader(tc.input)

			err := r.Start(context.Background(), "beta", tc.yes)
			s, _ := LoadSession(r.SessionPath)

			if tc.wantSwitch {
				if err != nil || s.Lab != "beta" {
					t.Fatalf("err = %v, session = %+v", err, s)
				}
				if fd.calls[0] != "rm breaklab-alpha" {
					t.Errorf("old container not removed first: %v", fd.calls)
				}
			} else {
				if err == nil || s.Lab != "alpha" || len(fd.calls) != 0 {
					t.Fatalf("err = %v, session = %+v, calls = %v", err, s, fd.calls)
				}
			}
		})
	}
}

func TestCheck(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()

	if _, err := r.Check(ctx); !errors.Is(err, ErrNoActiveLab) {
		t.Fatalf("no active lab: err = %v", err)
	}

	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	var gotEnv []string
	r.HostExec = func(_ context.Context, _, _ string, _, env []string) ([]byte, error) {
		gotEnv = env
		return []byte("PASS: fine\n"), nil
	}
	res, err := r.Check(ctx)
	if err != nil || !res.Passed || res.Message != "PASS: fine" || res.Variant != "first-variant" {
		t.Fatalf("pass: %+v, %v", res, err)
	}
	wantEnv := []string{"BREAKLAB_CONTAINER=breaklab-alpha", "BREAKLAB_PORT_80=8080"}
	if !reflect.DeepEqual(gotEnv, wantEnv) {
		t.Errorf("env = %v, want %v", gotEnv, wantEnv)
	}

	failErr := exec.Command("bash", "-c", "exit 1").Run()
	r.HostExec = func(context.Context, string, string, []string, []string) ([]byte, error) {
		return []byte("FAIL: down\n"), failErr
	}
	res, err = r.Check(ctx)
	if err != nil || res.Passed || res.Message != "FAIL: down" || res.Variant != "" {
		t.Fatalf("fail: %+v, %v", res, err)
	}
}

func TestCheckMissingScript(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(r.LabsDir, "alpha", "check.sh")); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Check(ctx); err == nil || !strings.Contains(err.Error(), "check.sh") {
		t.Fatalf("err = %v", err)
	}
}

// multiCheckLabYAML is a lab using the plural Checks field instead of
// the legacy singular Check, to exercise the multi-check path of
// Runner.Check.
func multiCheckLabYAML(id string) string {
	return "id: " + id + `
title: Test lab
stage: linux
difficulty: 1
estimated_minutes: 5
environment:
  dockerfile: Dockerfile
  ports: []
checks:
  - id: first-check
    script: checks/first.sh
    run_in: host
  - id: second-check
    script: checks/second.sh
    run_in: host
variants:
  - id: only-variant
    script: variants/only.sh
    weight: 1
`
}

func writeMultiCheckLab(t *testing.T, root, id string) {
	t.Helper()
	dir := filepath.Join(root, "labs", id)
	if err := os.MkdirAll(filepath.Join(dir, "checks"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "variants"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"lab.yaml":         multiCheckLabYAML(id),
		"Dockerfile":       "FROM scratch\n",
		"TICKET.md":        "# INC: " + id + " is down\n",
		"checks/first.sh":  "#!/usr/bin/env bash\necho PASS: first\n",
		"checks/second.sh": "#!/usr/bin/env bash\necho PASS: second\n",
		"variants/only.sh": "#!/usr/bin/env bash\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCheckRunsEveryNamedCheckAndAggregates(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	writeMultiCheckLab(t, filepath.Dir(r.LabsDir), "multi")
	ctx := context.Background()
	t.Setenv(forceVariantEnv, "only-variant")

	if err := r.Start(ctx, "multi", false); err != nil {
		t.Fatal(err)
	}

	var calledScripts []string
	r.HostExec = func(_ context.Context, _, _ string, args, _ []string) ([]byte, error) {
		script := args[0]
		calledScripts = append(calledScripts, filepath.Base(script))
		if strings.Contains(script, "second.sh") {
			return []byte("FAIL: second is broken\n"), exec.Command("bash", "-c", "exit 1").Run()
		}
		return []byte("PASS: first is fine\n"), nil
	}

	res, err := r.Check(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed {
		t.Error("overall Passed = true, want false: second-check failed")
	}
	if len(res.Results) != 2 {
		t.Fatalf("Results = %+v, want 2 entries", res.Results)
	}
	if res.Results[0].ID != "first-check" || !res.Results[0].Passed || res.Results[0].Message != "PASS: first is fine" {
		t.Errorf("Results[0] = %+v", res.Results[0])
	}
	if res.Results[1].ID != "second-check" || res.Results[1].Passed || res.Results[1].Message != "FAIL: second is broken" {
		t.Errorf("Results[1] = %+v", res.Results[1])
	}
	if !reflect.DeepEqual(calledScripts, []string{"first.sh", "second.sh"}) {
		t.Errorf("calledScripts = %v, want both run in order", calledScripts)
	}

	// Fix the second check and confirm overall Passed flips to true and
	// the attempt is marked solved.
	r.HostExec = func(context.Context, string, string, []string, []string) ([]byte, error) {
		return []byte("PASS: fine\n"), nil
	}
	res, err = r.Check(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Error("overall Passed = false after fixing both checks, want true")
	}
}

func TestStartMissingVariantScript(t *testing.T) {
	r, fd, _ := setup(t, "alpha")
	if err := os.Remove(filepath.Join(r.LabsDir, "alpha", "variants", "first.sh")); err != nil {
		t.Fatal(err)
	}
	if err := r.Start(context.Background(), "alpha", false); err == nil || !strings.Contains(err.Error(), "first.sh") {
		t.Fatalf("err = %v", err)
	}
	if len(fd.calls) != 0 {
		t.Errorf("docker touched before script check: %v", fd.calls)
	}
}

// noVariantLabYAML is a lab with no variants at all — the curriculum
// doc's standalone-lab model (Guided/Practice/single-state Break-Fix),
// which doesn't use chaos-injection.
func noVariantLabYAML(id string) string {
	return "id: " + id + `
title: Test lab, no variants
stage: linux
difficulty: 1
estimated_minutes: 5
environment:
  dockerfile: Dockerfile
  ports: []
check:
  script: check.sh
  run_in: host
`
}

func writeNoVariantLab(t *testing.T, root, id string) {
	t.Helper()
	dir := filepath.Join(root, "labs", id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"lab.yaml":   noVariantLabYAML(id),
		"Dockerfile": "FROM scratch\n",
		"TICKET.md":  "# INC: " + id + " is down\n",
		"check.sh":   "#!/usr/bin/env bash\necho PASS: fine\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestStartCheckResetStopWithNoVariants covers the full lifecycle for a
// lab that declares no variants at all: Start must not error with "no
// variants to choose from", launch must not try to copy/exec a
// nonexistent variant script, Check/Reset/Stop must all still work.
func TestStartCheckResetStopWithNoVariants(t *testing.T) {
	r, fd, _ := setup(t)
	writeNoVariantLab(t, filepath.Dir(r.LabsDir), "solo")
	ctx := context.Background()

	if err := r.Start(ctx, "solo", false); err != nil {
		t.Fatal(err)
	}
	for _, c := range fd.calls {
		if strings.Contains(c, containerVariantPath) {
			t.Errorf("a no-variant lab should never touch the variant path, got call %q", c)
		}
	}

	s, err := LoadSession(r.SessionPath)
	if err != nil || s == nil || s.Variant != "" {
		t.Fatalf("session = %+v, %v, want Variant empty", s, err)
	}

	res, err := r.Check(ctx)
	if err != nil || !res.Passed {
		t.Fatalf("Check = %+v, %v", res, err)
	}

	if err := r.Reset(ctx); err != nil {
		t.Fatal(err)
	}
	if err := r.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestResetKeepsVariantAndStopClears(t *testing.T) {
	r, fd, out := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	// Pretend the session picked the second variant; reset must keep it.
	s, _ := LoadSession(r.SessionPath)
	s.Variant = "second-variant"
	if err := SaveSession(r.SessionPath, *s); err != nil {
		t.Fatal(err)
	}

	fd.calls = nil
	out.Reset()
	if err := r.Reset(ctx); err != nil {
		t.Fatal(err)
	}
	if fd.calls[0] != "rm breaklab-alpha" || fd.calls[1] != "build breaklab-alpha" {
		t.Errorf("reset calls = %v", fd.calls)
	}
	s, _ = LoadSession(r.SessionPath)
	if s.Variant != "second-variant" {
		t.Errorf("variant after reset = %q", s.Variant)
	}
	if strings.Contains(out.String(), "variant") {
		t.Errorf("reset printed variant info: %q", out.String())
	}
	// Reset must not touch history.json: it still reflects Start's pick
	// ("first-variant"), not the session's manually-overwritten variant.
	if last, err := LastVariant(r.HistoryPath, "alpha"); err != nil || last != "first-variant" {
		t.Errorf("history after reset = %q, err %v, want unchanged %q", last, err, "first-variant")
	}

	if err := r.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadSession(r.SessionPath); s != nil {
		t.Errorf("session after stop = %+v", s)
	}
	if _, err := r.Check(ctx); !errors.Is(err, ErrNoActiveLab) {
		t.Errorf("check after stop: %v", err)
	}
}

func TestSession(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".breaklab", "session.json")

	if s, err := LoadSession(path); s != nil || err != nil {
		t.Fatalf("missing file: %v, %v", s, err)
	}

	in := Session{Lab: "alpha", Container: "breaklab-alpha", Variant: "v", HintsUsed: 3, AttemptIndex: 2}
	if err := SaveSession(path, in); err != nil {
		t.Fatal(err)
	}
	out, err := LoadSession(path)
	if err != nil || out.Lab != in.Lab || out.Variant != in.Variant || out.HintsUsed != in.HintsUsed || out.AttemptIndex != in.AttemptIndex {
		t.Fatalf("round trip: %+v, %v", out, err)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600", info.Mode().Perm())
	}

	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadSession(path); err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("corrupt file error = %v", err)
	}
}

func TestStartForceVariantSelectsExactly(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	t.Setenv(forceVariantEnv, "second-variant")

	if err := r.Start(context.Background(), "alpha", false); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSession(r.SessionPath)
	if err != nil || s == nil || s.Variant != "second-variant" {
		t.Fatalf("session = %+v, err %v", s, err)
	}
}

func TestStartForceVariantUnknownErrorsWithoutTouchingDocker(t *testing.T) {
	r, fd, _ := setup(t, "alpha")
	t.Setenv(forceVariantEnv, "no-such-variant")

	err := r.Start(context.Background(), "alpha", false)
	if err == nil || !strings.Contains(err.Error(), "no-such-variant") {
		t.Fatalf("err = %v", err)
	}
	if len(fd.calls) != 0 {
		t.Errorf("docker used despite invalid forced variant: %v", fd.calls)
	}
}

func TestStartExcludesPreviousVariantFromHistory(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	t.Setenv(forceVariantEnv, "")
	if err := RecordVariant(r.HistoryPath, "alpha", "first-variant"); err != nil {
		t.Fatal(err)
	}

	// With only two variants, each start must flip to the other one: the
	// previous start's pick becomes the next start's excluded variant.
	previous := "first-variant"
	for i := 0; i < 50; i++ {
		if err := r.Start(context.Background(), "alpha", true); err != nil {
			t.Fatal(err)
		}
		s, err := LoadSession(r.SessionPath)
		if err != nil || s.Variant == previous {
			t.Fatalf("iteration %d: session = %+v, err %v, previous %q", i, s, err, previous)
		}
		previous = s.Variant
	}
}

func TestStartSeedIsReproducible(t *testing.T) {
	t.Setenv(forceVariantEnv, "")
	seed := int64(7)

	pick := func() string {
		r, _, _ := setup(t, "alpha")
		t.Setenv(forceVariantEnv, "")
		r.Seed = &seed
		if err := r.Start(context.Background(), "alpha", false); err != nil {
			t.Fatal(err)
		}
		s, err := LoadSession(r.SessionPath)
		if err != nil {
			t.Fatal(err)
		}
		return s.Variant
	}

	first := pick()
	second := pick()
	if first != second {
		t.Errorf("seed %d picked %q then %q", seed, first, second)
	}
}

func TestStartRecordsVariantInHistory(t *testing.T) {
	r, _, _ := setup(t, "alpha")

	if err := r.Start(context.Background(), "alpha", false); err != nil {
		t.Fatal(err)
	}
	last, err := LastVariant(r.HistoryPath, "alpha")
	if err != nil || last != "first-variant" {
		t.Fatalf("last = %q, err %v", last, err)
	}
}

func TestCheckEnv(t *testing.T) {
	env, err := checkEnv("c", []string{"8080:80"})
	if err != nil || !reflect.DeepEqual(env, []string{"BREAKLAB_CONTAINER=c", "BREAKLAB_PORT_80=8080"}) {
		t.Errorf("env = %v, %v", env, err)
	}
	for _, bad := range []string{"80", "a:80", "8080:", "0:80", "8080:80:1", "1e3:80", "70000:80"} {
		if _, err := checkEnv("c", []string{bad}); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

func TestHostCheckWithRelativeLabsDir(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}

	// Use relative paths and the real host executor, as the CLI does.
	root := filepath.Dir(r.LabsDir)
	t.Chdir(root)
	r.LabsDir, r.RoadmapPath, r.SessionPath = "labs", "roadmap.yaml", filepath.Join(".breaklab", "session.json")
	r.HostExec = realHostExec

	res, err := r.Check(ctx)
	if err != nil || !res.Passed || res.Message != "PASS: fine" {
		t.Fatalf("check = %+v, %v", res, err)
	}
}

func TestContainerCheck(t *testing.T) {
	r, fd, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	labYAMLPath := filepath.Join(r.LabsDir, "alpha", "lab.yaml")
	if err := os.WriteFile(labYAMLPath, []byte(labYAML("alpha", "container")), 0o644); err != nil {
		t.Fatal(err)
	}

	fd.calls = nil
	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"cp breaklab-alpha " + containerCheckPath,
		"exec breaklab-alpha env BREAKLAB_CONTAINER=breaklab-alpha BREAKLAB_PORT_80=8080 bash " + containerCheckPath,
	}
	if !reflect.DeepEqual(fd.calls, want) {
		t.Errorf("calls = %v, want %v", fd.calls, want)
	}
}
