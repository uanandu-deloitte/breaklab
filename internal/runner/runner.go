// Package runner orchestrates start, check, reset and stop for a lab on top
// of the catalog and the docker wrapper.
package runner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/variants"
)

const (
	containerVariantPath = "/tmp/breaklab-variant.sh"
	containerCheckPath   = "/tmp/breaklab-check.sh"
	checkTimeout         = 30 * time.Second
)

// ErrNoActiveLab is returned by commands that need a started lab.
var ErrNoActiveLab = errors.New("no active lab; run `breaklab start <lab>` first")

// Docker is the subset of internal/docker the runner needs.
type Docker interface {
	Ping(ctx context.Context) error
	Build(ctx context.Context, image, contextDir, dockerfile string) error
	Run(ctx context.Context, name, image string, ports []string, privileged bool, volumes []string, pidsLimit int) error
	Exec(ctx context.Context, container string, argv ...string) ([]byte, error)
	Copy(ctx context.Context, src, container, dst string) error
	Remove(ctx context.Context, name string) error
	VolumeRemove(ctx context.Context, name string) error
	VolumeCreate(ctx context.Context, name string) error
}

// HostExecFunc runs a command on the host and returns combined output.
type HostExecFunc func(ctx context.Context, dir, name string, args, env []string) ([]byte, error)

// Runner holds everything the lab commands need.
type Runner struct {
	Docker       Docker
	LabsDir      string
	RoadmapPath  string
	SessionPath  string
	HistoryPath  string
	ProgressPath string

	// TemplatePath and PostmortemsDir locate templates/postmortem.md and
	// the postmortems/ directory for the Postmortem/ValidatePostmortem
	// methods.
	TemplatePath   string
	PostmortemsDir string

	HostExec    HostExecFunc
	OpenEditor  OpenEditorFunc
	In          io.Reader
	Out         io.Writer
	Interactive bool // stdin is a terminal, so prompting is possible

	// Seed, when non-nil, makes the next Start's variant selection
	// deterministic for that seed (see internal/variants.Pick).
	Seed *int64
}

// OpenEditorFunc opens path in an interactive editor and waits for it to
// exit.
type OpenEditorFunc func(ctx context.Context, path string) error

// New returns a Runner using the real host and the default repo layout.
func New(d Docker) *Runner {
	return &Runner{
		Docker:         d,
		LabsDir:        "labs",
		RoadmapPath:    "roadmap.yaml",
		SessionPath:    filepath.Join(".breaklab", "session.json"),
		HistoryPath:    filepath.Join(".breaklab", "history.json"),
		ProgressPath:   filepath.Join("progress", "progress.json"),
		TemplatePath:   filepath.Join("templates", "postmortem.md"),
		PostmortemsDir: "postmortems",
		HostExec:       realHostExec,
		OpenEditor:     realOpenEditor,
		In:             os.Stdin,
		Out:            os.Stdout,
	}
}

func realHostExec(ctx context.Context, dir, name string, args, env []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	return cmd.CombinedOutput()
}

// realOpenEditor opens path in $EDITOR (falling back to nano), split
// into argv with strings.Fields rather than passed to a shell, with
// stdio inherited so the editor is interactive.
func realOpenEditor(ctx context.Context, path string) error {
	argv := strings.Fields(os.Getenv("EDITOR"))
	if len(argv) == 0 {
		argv = []string{"nano"}
	}
	cmd := exec.CommandContext(ctx, argv[0], append(argv[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CheckResult is the outcome of `breaklab check`.
type CheckResult struct {
	Passed  bool
	Message string
	Variant string // set by the caller's output only when Passed

	// Results holds one entry per check when the lab verifies itself
	// through the plural Checks field instead of the legacy singular
	// Check — nil for every lab using the legacy field, same as before
	// Checks existed.
	Results []NamedCheckResult
}

// NamedCheckResult is one check's outcome within a multi-check lab's
// CheckResult.Results.
type NamedCheckResult struct {
	ID      string
	Passed  bool
	Message string
}

func containerName(labID string) string { return "breaklab-" + labID }

// forceVariantEnv names the environment variable that overrides variant
// selection entirely, for tests.
const forceVariantEnv = "BREAKLAB_FORCE_VARIANT"

// selectVariant returns the variant to apply for lab, in order of
// precedence: BREAKLAB_FORCE_VARIANT, then weighted random selection
// (seeded by r.Seed when set), excluding the lab's previously recorded
// variant per the no-immediate-repeat rule. A lab with no variants at
// all (the curriculum doc's standalone-lab model, which doesn't use
// chaos-injection) has nothing to select — the zero Variant, with an
// empty Script, which launch treats as "nothing to apply."
func (r *Runner) selectVariant(lab catalog.Lab) (catalog.Variant, error) {
	if len(lab.Variants) == 0 {
		return catalog.Variant{}, nil
	}

	if forced := os.Getenv(forceVariantEnv); forced != "" {
		return variants.Force(lab.Variants, forced)
	}

	previous, err := LastVariant(r.HistoryPath, lab.ID)
	if err != nil {
		return catalog.Variant{}, err
	}
	return variants.Pick(lab.Variants, previous, r.Seed)
}

// finalizeAttempt ends session's attempt record with the given status.
// It is idempotent (progress.MarkAbandoned no-ops on an attempt that
// isn't "active" any more — design.md Decision 6), so it's safe to call
// on a session whose attempt may already be solved.
func (r *Runner) finalizeAttempt(session *Session, status string) error {
	switch status {
	case progress.StatusAbandoned:
		return progress.MarkAbandoned(r.ProgressPath, session.AttemptIndex, session.HintsUsed)
	default:
		return fmt.Errorf("finalizeAttempt: unsupported status %q", status)
	}
}

func (r *Runner) findLab(id string) (catalog.Lab, error) {
	cat, err := catalog.LoadCatalog(r.LabsDir, r.RoadmapPath)
	if err != nil {
		return catalog.Lab{}, err
	}
	return findLabInCatalog(cat, id)
}

// findLabInCatalog looks up id in an already-loaded catalog, so a caller
// that needs the catalog for more than one purpose (e.g. Start, which
// also checks the capstone gate) doesn't reload and reparse every
// lab.yaml from disk a second time for the same request.
func findLabInCatalog(cat *catalog.Catalog, id string) (catalog.Lab, error) {
	for _, g := range cat.Stages {
		for _, l := range g.Labs {
			if l.ID == id {
				return l, nil
			}
		}
	}
	for _, w := range cat.Warnings {
		if strings.HasPrefix(w, id+":") {
			return catalog.Lab{}, fmt.Errorf("lab %q failed to load: %s", id, strings.TrimPrefix(w, id+": "))
		}
	}
	return catalog.Lab{}, fmt.Errorf("unknown lab %q; run `breaklab list` to see available labs", id)
}

// checkCapstoneUnlocked refuses to start lab if it belongs to the
// capstone stage and that stage isn't unlocked yet (every regular stage
// complete) — the capstone is an earned finale, enforced here so
// neither the CLI nor the web UI's action routes can bypass it. cat is
// the caller's already-loaded catalog, not reloaded here.
func (r *Runner) checkCapstoneUnlocked(cat *catalog.Catalog, lab catalog.Lab) error {
	capstone := cat.CapstoneStage()
	if capstone == nil || capstone.Stage.ID != lab.Stage {
		return nil
	}
	pf, err := progress.Load(r.ProgressPath)
	if err != nil {
		return err
	}
	status := progress.Capstone(cat, pf)
	if status != nil && !status.Unlocked {
		return fmt.Errorf("the capstone is locked until every other stage is complete; run `breaklab status` to see what's left")
	}
	return nil
}

// Start builds and starts a fresh container for the lab and applies a variant.
func (r *Runner) Start(ctx context.Context, labID string, yes bool) error {
	if err := r.Docker.Ping(ctx); err != nil {
		return err
	}

	cat, err := catalog.LoadCatalog(r.LabsDir, r.RoadmapPath)
	if err != nil {
		return err
	}
	lab, err := findLabInCatalog(cat, labID)
	if err != nil {
		return err
	}
	if err := r.checkCapstoneUnlocked(cat, lab); err != nil {
		return err
	}

	session, err := LoadSession(r.SessionPath)
	if err != nil {
		return err
	}
	if session != nil && session.Lab != lab.ID {
		if !yes && !r.confirm(fmt.Sprintf("Lab %q is active. Stop it and start %q?", session.Lab, lab.ID)) {
			return fmt.Errorf("lab %q is still active; run `breaklab stop` first or pass --yes", session.Lab)
		}
		if err := r.Docker.Remove(ctx, session.Container); err != nil {
			return err
		}
		if err := r.finalizeAttempt(session, progress.StatusAbandoned); err != nil {
			return err
		}
		if err := ClearSession(r.SessionPath); err != nil {
			return err
		}
	} else if session != nil {
		// Restarting the same lab without `reset` also begins a new
		// attempt (a new variant may be selected below); finalize the
		// one being replaced first, same as switching labs.
		if err := r.finalizeAttempt(session, progress.StatusAbandoned); err != nil {
			return err
		}
	}

	variant, err := r.selectVariant(lab)
	if err != nil {
		return err
	}
	index, err := progress.AppendActive(r.ProgressPath, lab.ID, variant.ID)
	if err != nil {
		return err
	}
	if err := r.launch(ctx, lab, variant, index); err != nil {
		return err
	}
	if err := RecordVariant(r.HistoryPath, lab.ID, variant.ID); err != nil {
		return err
	}

	ticket, err := catalog.ReadTicket(r.LabsDir, lab.ID)
	if err != nil {
		return fmt.Errorf("lab started but its ticket is unreadable: %w", err)
	}
	fmt.Fprintf(r.Out, "%s\n", strings.TrimRight(string(ticket), "\n"))
	shell := fmt.Sprintf("docker exec -it %s bash", containerName(lab.ID))
	if !strings.Contains(string(ticket), shell) {
		fmt.Fprintf(r.Out, "\nGet a shell: %s\n", shell)
	}
	return nil
}

// launch removes any old container, builds, runs, applies the variant and
// writes the session. On failure the half-started container is removed.
// attemptIndex is written onto the session as-is: Start passes a freshly
// appended progress.json index, Reset passes the existing session's
// index through unchanged, so a reset never starts a new attempt.
func (r *Runner) launch(ctx context.Context, lab catalog.Lab, variant catalog.Variant, attemptIndex int) (err error) {
	name := containerName(lab.ID)
	labDir := filepath.Join(r.LabsDir, lab.ID)
	variantPath := filepath.Join(labDir, variant.Script)

	// A lab with no variants (catalog.Lab.Variants empty) has nothing to
	// apply — its starting state is whatever the Dockerfile itself
	// builds, with no separate breakage step layered on top.
	hasVariant := variant.Script != ""
	if hasVariant {
		if _, statErr := os.Stat(variantPath); statErr != nil {
			return fmt.Errorf("lab %q: variant script %s not found", lab.ID, variantPath)
		}
	}

	if _, err := checkEnv(name, lab.Environment.Ports); err != nil {
		return err
	}

	// The old container (if any) must go first: a Named volume can't be
	// removed and recreated while a container, even a stopped one,
	// still references it.
	if err := r.Docker.Remove(ctx, name); err != nil {
		return err
	}

	volumes, err := r.resolveVolumes(ctx, lab)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = r.Docker.Remove(context.Background(), name)
		}
	}()

	dockerfile := filepath.Join(labDir, lab.Environment.Dockerfile)
	if err := r.Docker.Build(ctx, name, labDir, dockerfile); err != nil {
		return err
	}
	if err := r.Docker.Run(ctx, name, name, lab.Environment.Ports, lab.Environment.Privileged, volumes, lab.Environment.PidsLimit); err != nil {
		return err
	}
	if hasVariant {
		if err := r.Docker.Copy(ctx, variantPath, name, containerVariantPath); err != nil {
			return err
		}
		if out, err := r.Docker.Exec(ctx, name, "bash", containerVariantPath); err != nil {
			return fmt.Errorf("applying the lab's breakage failed: %w\n%s", err, strings.TrimSpace(string(out)))
		}
		// Best-effort: a leftover temp file in a throwaway container that
		// gets rebuilt on every Start/Reset is harmless, and a handful of
		// labs leave the container briefly too short on some resource
		// (e.g. pids) for even this tiny exec to land right away.
		_, _ = r.Docker.Exec(ctx, name, "rm", "-f", containerVariantPath)
	}

	return SaveSession(r.SessionPath, Session{
		Lab:          lab.ID,
		Container:    name,
		Variant:      variant.ID,
		AttemptIndex: attemptIndex,
	})
}

// resolveVolumes prepares a fresh mount for each of lab's declared
// volumes and returns the "source:containerPath" strings Run expects.
// A plain volume becomes an emptied host directory under a
// session-scoped "volumes" directory next to SessionPath. A Named
// volume becomes a daemon-managed Docker volume instead (recreated from
// scratch each time), for labs that need real multi-UID permission
// enforcement a host bind mount doesn't reliably give on every backend.
// Either way, a lab never inherits a previous attempt's leftover state —
// consistent with launch() always rebuilding the container from scratch.
func (r *Runner) resolveVolumes(ctx context.Context, lab catalog.Lab) ([]string, error) {
	if len(lab.Environment.Volumes) == 0 {
		return nil, nil
	}
	base := filepath.Join(filepath.Dir(r.SessionPath), "volumes", lab.ID)
	mounts := make([]string, 0, len(lab.Environment.Volumes))
	for _, v := range lab.Environment.Volumes {
		if !catalog.ValidID(v.Name) {
			return nil, fmt.Errorf("lab %q: invalid volume name %q", lab.ID, v.Name)
		}
		if v.Named {
			volName := "breaklab-" + lab.ID + "-" + v.Name
			if err := r.Docker.VolumeRemove(ctx, volName); err != nil {
				return nil, fmt.Errorf("clearing volume %q: %w", v.Name, err)
			}
			if err := r.Docker.VolumeCreate(ctx, volName); err != nil {
				return nil, fmt.Errorf("creating volume %q: %w", v.Name, err)
			}
			mounts = append(mounts, volName+":"+v.ContainerPath)
			continue
		}
		dir := filepath.Join(base, v.Name)
		if err := os.RemoveAll(dir); err != nil {
			return nil, fmt.Errorf("clearing volume %q: %w", v.Name, err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("creating volume %q: %w", v.Name, err)
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		mounts = append(mounts, abs+":"+v.ContainerPath)
	}
	return mounts, nil
}

// Check runs the active lab's check script (or, for a lab using the
// plural Checks field, every one of them in order).
func (r *Runner) Check(ctx context.Context) (CheckResult, error) {
	session, err := LoadSession(r.SessionPath)
	if err != nil {
		return CheckResult{}, err
	}
	if session == nil {
		return CheckResult{}, ErrNoActiveLab
	}
	if err := r.Docker.Ping(ctx); err != nil {
		return CheckResult{}, err
	}

	lab, err := r.findLab(session.Lab)
	if err != nil {
		return CheckResult{}, err
	}

	labDir := filepath.Join(r.LabsDir, lab.ID)
	env, err := checkEnv(session.Container, lab.Environment.Ports)
	if err != nil {
		return CheckResult{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	var res CheckResult
	if len(lab.Checks) > 0 {
		res.Passed = true
		var messages []string
		for _, c := range lab.Checks {
			passed, message, err := r.runOneCheck(ctx, labDir, session.Container, c.Script, c.RunIn, env)
			if err != nil {
				return CheckResult{}, fmt.Errorf("check %q: %w", c.ID, err)
			}
			res.Results = append(res.Results, NamedCheckResult{ID: c.ID, Passed: passed, Message: message})
			messages = append(messages, c.ID+": "+message)
			if !passed {
				res.Passed = false
			}
		}
		res.Message = strings.Join(messages, "; ")
	} else {
		res.Passed, res.Message, err = r.runOneCheck(ctx, labDir, session.Container, lab.Check.Script, lab.Check.RunIn, env)
		if err != nil {
			return CheckResult{}, err
		}
	}

	if res.Passed {
		res.Variant = session.Variant
		if err := progress.MarkSolved(r.ProgressPath, session.AttemptIndex, session.HintsUsed); err != nil {
			return CheckResult{}, err
		}
	} else {
		if err := progress.IncrementFailedChecks(r.ProgressPath, session.AttemptIndex); err != nil {
			return CheckResult{}, err
		}
	}
	return res, nil
}

// runOneCheck runs a single check script — either the legacy Check or
// one entry of Checks — against the active session's container/host and
// returns whether it passed and its last printed line as the verdict
// message. ctx must already carry Check's overall checkTimeout; every
// check in a multi-check lab shares that one budget rather than each
// getting its own.
func (r *Runner) runOneCheck(ctx context.Context, labDir, container, scriptRel, runIn string, env []string) (passed bool, message string, err error) {
	script := filepath.Join(labDir, scriptRel)
	if _, statErr := os.Stat(script); statErr != nil {
		return false, "", fmt.Errorf("check script %s not found", script)
	}
	// The host check runs with the lab directory as its working directory, so
	// a relative script path would resolve against the wrong place.
	if script, err = filepath.Abs(script); err != nil {
		return false, "", err
	}

	var out []byte
	switch runIn {
	case "container":
		if err := r.Docker.Copy(ctx, script, container, containerCheckPath); err != nil {
			return false, "", err
		}
		argv := append(append([]string{"env"}, env...), "bash", containerCheckPath)
		out, err = r.Docker.Exec(ctx, container, argv...)
	default:
		out, err = r.HostExec(ctx, labDir, "bash", []string{script}, env)
	}

	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		if ctx.Err() != nil {
			return false, "", fmt.Errorf("check timed out after %s", checkTimeout)
		}
		return false, "", fmt.Errorf("running check: %w", err)
	}

	passed = err == nil
	message = lastLine(out)
	if message == "" {
		message = map[bool]string{true: "check passed", false: "check failed"}[passed]
	}
	return passed, message, nil
}

// Reset recreates the active lab's container with the same variant.
func (r *Runner) Reset(ctx context.Context) error {
	session, err := LoadSession(r.SessionPath)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrNoActiveLab
	}
	if err := r.Docker.Ping(ctx); err != nil {
		return err
	}

	lab, err := r.findLab(session.Lab)
	if err != nil {
		return err
	}
	var variant catalog.Variant
	// A lab with no variants always has an empty session.Variant (see
	// selectVariant) — nothing to look up, and launch treats the zero
	// Variant as "nothing to apply," same as the original Start did.
	found := len(lab.Variants) == 0 && session.Variant == ""
	for _, v := range lab.Variants {
		if v.ID == session.Variant {
			variant, found = v, true
		}
	}
	if !found {
		return fmt.Errorf("variant of the active session no longer exists in lab %q; run `breaklab start %s`", lab.ID, lab.ID)
	}
	// Reset continues the same attempt: same AttemptIndex, so
	// progress.json's started_at and failed_checks span every reset.
	return r.launch(ctx, lab, variant, session.AttemptIndex)
}

// Stop removes the active container and clears the session.
func (r *Runner) Stop(ctx context.Context) error {
	session, err := LoadSession(r.SessionPath)
	if err != nil {
		return err
	}
	if session == nil {
		fmt.Fprintln(r.Out, "No lab is active.")
		return nil
	}
	if err := r.Docker.Ping(ctx); err != nil {
		return err
	}
	if err := r.Docker.Remove(ctx, session.Container); err != nil {
		return err
	}
	if err := r.finalizeAttempt(session, progress.StatusAbandoned); err != nil {
		return err
	}
	return ClearSession(r.SessionPath)
}

func (r *Runner) confirm(prompt string) bool {
	if !r.Interactive {
		return false
	}
	fmt.Fprintf(r.Out, "%s [y/N] ", prompt)
	line, _ := bufio.NewReader(r.In).ReadString('\n')
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes"
}

// checkEnv builds BREAKLAB_CONTAINER and BREAKLAB_PORT_<container> entries
// from "host:container" port mappings.
func checkEnv(container string, ports []string) ([]string, error) {
	env := []string{"BREAKLAB_CONTAINER=" + container}
	for _, p := range ports {
		host, ctr, ok := strings.Cut(p, ":")
		if !ok || !validPort(host) || !validPort(ctr) {
			return nil, fmt.Errorf("invalid port mapping %q: want host:container, e.g. 8080:80", p)
		}
		env = append(env, fmt.Sprintf("BREAKLAB_PORT_%s=%s", ctr, host))
	}
	return env, nil
}

func validPort(s string) bool {
	n, err := strconv.Atoi(s)
	return err == nil && n >= 1 && n <= 65535 && strconv.Itoa(n) == s
}

// lastLine returns the last non-empty line: a check script prints its verdict
// last, after any noise (such as curl's stderr) that came before it.
func lastLine(out []byte) string {
	lines := strings.Split(string(out), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if line := strings.TrimSpace(lines[i]); line != "" {
			return line
		}
	}
	return ""
}
