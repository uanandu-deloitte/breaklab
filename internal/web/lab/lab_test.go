package lab

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
)

func fixtureLabsDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "web-down"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "web-down", "TICKET.md"), []byte("# INC-101: Website down\nRestore the homepage.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func fixtureCatalogForLabDetail() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{
				ID:    "web-down",
				Title: "The website returns nothing",
				Hints: []string{"SECRET_GENERAL_HINT"},
				Variants: []catalog.Variant{
					{ID: "config-typo", Hints: []string{"SECRET_VARIANT_HINT"}},
				},
			},
		}},
	}}
}

func TestBuildLabDetailNeverLeaksVariantOrHints(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{}

	view, found, err := BuildDetail(context.Background(), cat, pf, nil, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if strings.Contains(view.Ticket, "SECRET_GENERAL_HINT") || strings.Contains(view.Ticket, "SECRET_VARIANT_HINT") {
		t.Errorf("ticket leaks hint text: %q", view.Ticket)
	}
	if strings.Contains(view.Ticket, "config-typo") {
		t.Errorf("ticket leaks the variant id: %q", view.Ticket)
	}
	if view.HintsUsed != 0 {
		t.Errorf("HintsUsed = %d, want 0 for an untouched lab", view.HintsUsed)
	}
	if view.Command != "breaklab start web-down" {
		t.Errorf("Command = %q", view.Command)
	}
	if view.HintSlots != nil {
		t.Errorf("HintSlots = %+v, want nil for a non-active lab", view.HintSlots)
	}
}

func TestBuildLabDetailActiveLabHintSlotsNeverLeakText(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: time.Now(), Status: progress.StatusActive},
	}}
	session := &runner.Session{Lab: "web-down", Variant: "config-typo", AttemptIndex: 0, HintsUsed: 1}

	view, found, err := BuildDetail(context.Background(), cat, pf, session, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	// 1 lab hint + 1 variant hint = 2 total slots.
	if len(view.HintSlots) != 2 {
		t.Fatalf("HintSlots = %+v, want 2 slots", view.HintSlots)
	}
	if !view.HintSlots[0].Revealed {
		t.Errorf("HintSlots[0].Revealed = false, want true (hintsUsed=1)")
	}
	if view.HintSlots[1].Revealed {
		t.Errorf("HintSlots[1].Revealed = true, want false (hintsUsed=1)")
	}
	if view.HintSlots[1].Cost != hintCost(1) {
		t.Errorf("HintSlots[1].Cost = %d, want %d", view.HintSlots[1].Cost, hintCost(1))
	}
	if view.TicketID != webcore.TicketID(0) {
		t.Errorf("TicketID = %q, want %q", view.TicketID, webcore.TicketID(0))
	}
}

func TestBuildLabDetailDiagnosticLogDegradesOnError(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: time.Now(), Status: progress.StatusActive},
	}}
	session := &runner.Session{Lab: "web-down", Variant: "config-typo", Container: "breaklab-web-down", AttemptIndex: 0}

	failingLogs := webcore.LogsFunc(func(ctx context.Context, container string, tailLines int) ([]byte, error) {
		return nil, errors.New("docker unreachable")
	})
	view, found, err := BuildDetail(context.Background(), cat, pf, session, "web-down", labsDir, failingLogs, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if view.DiagnosticLog != "" {
		t.Errorf("DiagnosticLog = %q, want empty on a Docker error", view.DiagnosticLog)
	}

	workingLogs := webcore.LogsFunc(func(ctx context.Context, container string, tailLines int) ([]byte, error) {
		if container != "breaklab-web-down" {
			t.Errorf("logs called with container %q, want breaklab-web-down", container)
		}
		return []byte("line one\nline two\n"), nil
	})
	view, _, _ = BuildDetail(context.Background(), cat, pf, session, "web-down", labsDir, workingLogs, false)
	if view.DiagnosticLog != "line one\nline two" {
		t.Errorf("DiagnosticLog = %q", view.DiagnosticLog)
	}
}

func TestBuildLabDetailActiveLabShowsCheckCommand(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	started := time.Now().Add(-5 * time.Minute)
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: started, Status: progress.StatusActive},
	}}
	session := &runner.Session{Lab: "web-down", AttemptIndex: 0, HintsUsed: 3}

	view, found, err := BuildDetail(context.Background(), cat, pf, session, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if view.Command != "breaklab check" {
		t.Errorf("Command = %q, want breaklab check", view.Command)
	}
	if view.HintsUsed != 3 {
		t.Errorf("HintsUsed = %d, want 3 (live from session)", view.HintsUsed)
	}
}

func TestBuildLabDetailNonActiveLabUsesSnapshot(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: time.Now().Add(-time.Hour), Status: progress.StatusAbandoned, HintsUsed: 2},
	}}

	// No session at all: this lab's own history still shows its snapshot.
	view, found, err := BuildDetail(context.Background(), cat, pf, nil, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if view.Command != "breaklab start web-down" {
		t.Errorf("Command = %q, want breaklab start web-down", view.Command)
	}
	if view.HintsUsed != 2 {
		t.Errorf("HintsUsed = %d, want 2 (from the latest attempt)", view.HintsUsed)
	}
}

// fakeReadCloser blocks Read until Close (simulating a long-running
// "docker logs -f" stream) and records whether Close was ever called.
type fakeReadCloser struct {
	pr     *io.PipeReader
	pw     *io.PipeWriter
	once   sync.Once
	closed chan struct{}
}

func newFakeReadCloser() *fakeReadCloser {
	pr, pw := io.Pipe()
	return &fakeReadCloser{pr: pr, pw: pw, closed: make(chan struct{})}
}

func (f *fakeReadCloser) Read(p []byte) (int, error) { return f.pr.Read(p) }

func (f *fakeReadCloser) Close() error {
	f.once.Do(func() { close(f.closed) })
	return f.pr.Close()
}

func TestLabLogsStreamWrongLabIs404(t *testing.T) {
	dir := t.TempDir()
	if err := runner.SaveSession(filepath.Join(dir, "session.json"), runner.Session{
		Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo", AttemptIndex: 0,
	}); err != nil {
		t.Fatal(err)
	}
	cfg := webcore.Config{
		SessionPath:         filepath.Join(dir, "session.json"),
		ContainerLogsFollow: func(ctx context.Context, container string) (io.ReadCloser, error) { return newFakeReadCloser(), nil },
	}

	req := httptest.NewRequest("GET", "/labs/some-other-lab/logs/stream", nil)
	req.SetPathValue("id", "some-other-lab")
	rec := httptest.NewRecorder()
	labLogsStreamHandler(cfg)(rec, req)
	if rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

// fakeOneShotReadCloser hands back a fixed body once, then blocks (like
// a real "docker logs -f" would once it's caught up) until Close.
type fakeOneShotReadCloser struct {
	r      io.Reader
	block  chan struct{}
	closed chan struct{}
	once   sync.Once
}

func newFakeOneShotReadCloser(body string) *fakeOneShotReadCloser {
	return &fakeOneShotReadCloser{r: strings.NewReader(body), block: make(chan struct{}), closed: make(chan struct{})}
}

func (f *fakeOneShotReadCloser) Read(p []byte) (int, error) {
	n, err := f.r.Read(p)
	if err == io.EOF {
		<-f.block // then hang, like a real -f stream with no more output yet
	}
	return n, err
}

func (f *fakeOneShotReadCloser) Close() error {
	f.once.Do(func() {
		close(f.block)
		close(f.closed)
	})
	return nil
}

func TestLabLogsStreamEscapesLogContent(t *testing.T) {
	dir := t.TempDir()
	if err := runner.SaveSession(filepath.Join(dir, "session.json"), runner.Session{
		Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo", AttemptIndex: 0,
	}); err != nil {
		t.Fatal(err)
	}
	frc := newFakeOneShotReadCloser("<script>alert(1)</script>\n")
	cfg := webcore.Config{
		SessionPath:         filepath.Join(dir, "session.json"),
		ContainerLogsFollow: func(ctx context.Context, container string) (io.ReadCloser, error) { return frc, nil },
	}

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/labs/web-down/logs/stream", nil).WithContext(ctx)
	req.SetPathValue("id", "web-down")
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		labLogsStreamHandler(cfg)(rec, req)
		close(done)
	}()

	// Give the handler a moment to read and write the one line, then
	// disconnect so the handler returns.
	time.Sleep(200 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler never returned after disconnect")
	}

	body := rec.Body.String()
	if strings.Contains(body, "<script>") {
		t.Fatalf("raw <script> tag present in SSE body — log content must be HTML-escaped: %q", body)
	}
	if !strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatalf("expected escaped log line in SSE body, got: %q", body)
	}
	if !strings.Contains(body, `<span class="log-line">`) {
		t.Fatalf("expected each line wrapped in a log-line span, got: %q", body)
	}
}

func TestLabLogsStreamClosesOnClientDisconnect(t *testing.T) {
	dir := t.TempDir()
	if err := runner.SaveSession(filepath.Join(dir, "session.json"), runner.Session{
		Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo", AttemptIndex: 0,
	}); err != nil {
		t.Fatal(err)
	}
	frc := newFakeReadCloser()
	cfg := webcore.Config{
		SessionPath: filepath.Join(dir, "session.json"),
		ContainerLogsFollow: func(ctx context.Context, container string) (io.ReadCloser, error) {
			if container != "breaklab-web-down" {
				t.Errorf("container = %q, want breaklab-web-down", container)
			}
			return frc, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/labs/web-down/logs/stream", nil).WithContext(ctx)
	req.SetPathValue("id", "web-down")
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		labLogsStreamHandler(cfg)(rec, req)
		close(done)
	}()

	cancel()
	select {
	case <-frc.closed:
	case <-time.After(2 * time.Second):
		t.Fatal("ContainerLogsFollow's ReadCloser was never closed after client disconnect")
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler never returned after Close")
	}
}

func TestBuildLabDetailReadsHypothesisFromLatestAttempt(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{Attempts: []progress.Attempt{
		{Lab: "web-down", Variant: "config-typo", StartedAt: time.Now(), Status: progress.StatusActive, Hypothesis: "probably a config typo"},
	}}
	session := &runner.Session{Lab: "web-down", Variant: "config-typo", AttemptIndex: 0}

	view, found, err := BuildDetail(context.Background(), cat, pf, session, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if view.Hypothesis != "probably a config typo" {
		t.Errorf("Hypothesis = %q", view.Hypothesis)
	}

	// Non-active (snapshot) view reads the same way.
	view, found, err = BuildDetail(context.Background(), cat, pf, nil, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if view.Hypothesis != "probably a config typo" {
		t.Errorf("non-active Hypothesis = %q", view.Hypothesis)
	}
}

func TestBuildLabDetailPrefersReadmeOverTicket(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	if err := os.WriteFile(filepath.Join(labsDir, "web-down", "README.md"), []byte("# From README\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{}

	view, found, err := BuildDetail(context.Background(), cat, pf, nil, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if !strings.Contains(view.Ticket, "From README") {
		t.Errorf("Ticket = %q, want README.md content preferred over TICKET.md", view.Ticket)
	}
}

func TestBuildLabDetailRendersLessonAndType(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	if err := os.WriteFile(filepath.Join(labsDir, "web-down", "lesson.md"), []byte("# Concept\n\nSome **bold** text."), 0o644); err != nil {
		t.Fatal(err)
	}
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "web-down", Title: "The website returns nothing", Type: "guided", PathID: "linux-fundamentals", Module: "orientation", Objectives: []string{"collect operating context"}, Lesson: "lesson.md"},
		}},
	}}
	pf := &progress.ProgressFile{}

	view, found, err := BuildDetail(context.Background(), cat, pf, nil, "web-down", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if view.Type != "guided" || view.PathID != "linux-fundamentals" || view.Module != "orientation" {
		t.Errorf("Type/PathID/Module = %q/%q/%q", view.Type, view.PathID, view.Module)
	}
	if len(view.Objectives) != 1 || view.Objectives[0] != "collect operating context" {
		t.Errorf("Objectives = %v", view.Objectives)
	}
	if !strings.Contains(string(view.Lesson), "<h1>Concept</h1>") || !strings.Contains(string(view.Lesson), "<strong>bold</strong>") {
		t.Errorf("Lesson = %q, want rendered markdown", view.Lesson)
	}
}

func TestBuildLabDetailLockedLabCannotLaunch(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	if err := os.MkdirAll(filepath.Join(labsDir, "lnx-002"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(labsDir, "lnx-002", "TICKET.md"), []byte("ticket\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "web-down", Title: "The website returns nothing"},
			{ID: "lnx-002", Title: "Second lab", Prerequisites: []string{"web-down"}},
		}},
	}}
	pf := &progress.ProgressFile{}

	view, found, err := BuildDetail(context.Background(), cat, pf, nil, "lnx-002", labsDir, nil, false)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if view.Status != progress.LabStatusLocked {
		t.Fatalf("Status = %q, want locked", view.Status)
	}
	if view.CanLaunch {
		t.Error("CanLaunch = true for a locked lab, want false")
	}
	if len(view.LockedOnTitles) != 1 || view.LockedOnTitles[0] != "The website returns nothing" {
		t.Errorf("LockedOnTitles = %v", view.LockedOnTitles)
	}
}

func TestBuildLabDetailChallengeMode(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{}

	view, found, err := BuildDetail(context.Background(), cat, pf, nil, "web-down", labsDir, nil, true)
	if err != nil || !found {
		t.Fatalf("found = %v, err = %v", found, err)
	}
	if !view.ChallengeMode {
		t.Error("ChallengeMode = false, want true for ?mode=challenge")
	}
}

func TestBuildLabDetailUnknownIDNotFound(t *testing.T) {
	labsDir := fixtureLabsDir(t)
	cat := fixtureCatalogForLabDetail()
	pf := &progress.ProgressFile{}

	_, found, err := BuildDetail(context.Background(), cat, pf, nil, "no-such-lab", labsDir, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Error("expected found = false for an unknown lab id")
	}
}
