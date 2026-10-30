package lab

import (
	"bufio"
	"context"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"breaklab/internal/catalog"
	"breaklab/internal/markdown"
	"breaklab/internal/progress"
	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
)

// HintSlotView is one hint's lock state in the lab detail page's hint
// list. Cost is a cosmetic point value only — display label, never
// deducted from or added to any stored total — and is meaningless once
// Revealed (the revealed hint's own text comes from Flash.Hint, not from
// here: this package never carries an already-used hint's text back to
// the client beyond the one just revealed by the current request).
type HintSlotView struct {
	Revealed bool
	Cost     int
}

// LabDetailView is everything templates/web/lab.html renders. It never
// carries a variant id or any hint text beyond the one just revealed by
// this request — not used hints from earlier requests, solved or not —
// only a count (design.md Decision 4).
type LabDetailView struct {
	ID               string
	Title            string
	StageTitle       string
	Status           string
	Difficulty       int
	EstimatedMinutes int
	Tags             []string
	Ticket           string
	TicketID         string
	Ports            []webcore.PortNodeView
	HintsUsed        int
	HintSlots        []HintSlotView // set only when Active
	DiagnosticLog    string         // tail of the container's own logs; set only when Active and available
	Command          string
	Active           bool   // this lab is the current session: Check/Hint/Stop apply to it
	ShellCmd         string // set only when Active
	CanLaunch        bool   // no other lab's session is blocking a plain Start
	BlockingLab      string // set when a different lab's session blocks a plain Start

	// Summary and the Learn/Guided/Practice fields are this lab's loop
	// on-ramp content (catalog.Lab.Summary / .Loop), shown before the
	// Break/Fix ticket. All zero-valued for a lab with no loop content
	// yet — their sections just don't render.
	Summary  string
	Learn    string
	Guided   []GuidedStepView
	Practice PracticeBlockView

	// ChallengeMode is set from the ?mode=challenge query/form value and
	// carried through every action's redirect (see redirectToLab). It
	// never gates the Hint action server-side — it only hides the
	// Hints section and its reveal form from this page's own markup.
	ChallengeMode bool

	// Hypothesis is the learner's own saved guess at the cause for this
	// lab's latest attempt, if any — read straight from progress.json
	// (see buildLabDetail), not session state, since it can be written
	// and edited any time the attempt is still active.
	Hypothesis string

	// Type, PathID, Module, and Objectives are the curriculum doc's
	// standalone-lab metadata (catalog.Lab's new fields). All empty for
	// every lab written before this existed — the Type badge and
	// breadcrumb just don't render for them.
	Type       string
	PathID     string
	Module     string
	Objectives []string

	// Lesson is this lab's rendered markdown concept primer
	// (catalog.Lab.Lesson), distinct from the plain-text Learn field
	// above — a lab sets one or the other, never both. Empty for every
	// lab without a lesson.md.
	Lesson template.HTML

	// LockedOnTitles lists the titles of this lab's not-yet-complete
	// Prerequisites, set only when Status is "locked" — CanLaunch is
	// false and the Launch section explains what to finish first
	// instead of offering to start.
	LockedOnTitles []string
}

// GuidedStepView is one step of the lab's guided walkthrough.
type GuidedStepView struct {
	Do  string
	Why string
}

// PracticeBlockView is the lab's light rehearsal exercise.
type PracticeBlockView struct {
	Task      string
	Questions []PracticeQuestionView
}

// PracticeQuestionView is one revealable question/answer pair.
type PracticeQuestionView struct {
	Prompt string
	Answer string
}

// hintCost is the cosmetic point value shown next to a locked hint at
// position hintIndex (0-based) within a lab's combined hint list. It is
// a display-only label — nothing is ever deducted from or added to any
// stored total; there is no scoring system behind it.
func hintCost(hintIndex int) int {
	return 5 * (hintIndex + 1)
}

// BuildDetail assembles LabDetailView for labID. The returned bool
// is false when labID doesn't match any lab (the caller should 404);
// ticket reading uses the matched lab's own catalog-validated ID, never
// the raw request path value (design.md Decision 7's same principle).
// logs is called only with session.Container, never request input, and
// only while active — a nil logs or a lookup error just leaves
// DiagnosticLog empty rather than failing the page (same contract as
// topology's status func).
func BuildDetail(ctx context.Context, cat *catalog.Catalog, pf *progress.ProgressFile, session *runner.Session, labID, labsDir string, logs webcore.LogsFunc, challengeMode bool) (LabDetailView, bool, error) {
	progress.ApplyStatus(cat, pf)

	lab, stageTitle, found := cat.FindLab(labID)
	if !found {
		return LabDetailView{}, false, nil
	}

	ticket, err := catalog.ReadTicket(labsDir, lab.ID)
	if err != nil {
		return LabDetailView{}, false, err
	}

	var lesson template.HTML
	if lab.Lesson != "" {
		src, err := os.ReadFile(filepath.Join(labsDir, lab.ID, lab.Lesson))
		if err != nil {
			return LabDetailView{}, false, err
		}
		lesson = markdown.Render(string(src))
	}

	hintsUsed := 0
	command := "breaklab start " + lab.ID
	active := session != nil && session.Lab == lab.ID
	canLaunch := true
	blockingLab := ""
	ticketIdx, haveTicketIdx := -1, false
	var lockedOnTitles []string

	// Hypothesis is never session-tracked (unlike HintsUsed) — it's
	// written straight to progress.json by SetHypothesis — so it's read
	// from the latest attempt the same way whether or not this lab is
	// the active session.
	hypothesis := ""
	if attempt, ok := progress.LatestAttempt(pf, lab.ID); ok {
		hypothesis = attempt.Hypothesis
	}

	if active {
		hintsUsed = session.HintsUsed
		command = "breaklab check"
		ticketIdx, haveTicketIdx = session.AttemptIndex, true
	} else {
		if attempt, ok := progress.LatestAttempt(pf, lab.ID); ok {
			hintsUsed = attempt.HintsUsed
		}
		if idx, ok := progress.LatestAttemptIndex(pf, lab.ID); ok {
			ticketIdx, haveTicketIdx = idx, true
		}
		if session != nil {
			canLaunch = false
			blockingLab = session.Lab
		}
	}

	// A locked lab can never be launched, regardless of session state —
	// this overrides whatever the concurrent-session check above
	// decided, the same way a locked stage's capstone is refused
	// server-side in the runner, not just hidden in the UI.
	if lab.Status == progress.LabStatusLocked {
		canLaunch = false
		for _, id := range lab.Prerequisites {
			preLab, _, ok := cat.FindLab(id)
			if ok && preLab.Status != progress.LabStatusComplete {
				lockedOnTitles = append(lockedOnTitles, preLab.Title)
			}
		}
	}

	var ports []webcore.PortNodeView
	for _, p := range lab.Environment.Ports {
		host, container, ok := strings.Cut(p, ":")
		pv := webcore.PortNodeView{HostPort: host, ContainerPort: container}
		if ok && host != "" {
			pv.OpenHref = "http://localhost:" + host
		}
		ports = append(ports, pv)
	}

	view := LabDetailView{
		ID:               lab.ID,
		Title:            lab.Title,
		StageTitle:       stageTitle,
		Status:           lab.Status,
		Difficulty:       lab.Difficulty,
		EstimatedMinutes: lab.EstimatedMinutes,
		Tags:             lab.Tags,
		Ticket:           string(ticket),
		Ports:            ports,
		HintsUsed:        hintsUsed,
		Command:          command,
		Active:           active,
		CanLaunch:        canLaunch,
		BlockingLab:      blockingLab,
		Summary:          lab.Summary,
		Learn:            lab.Loop.Learn,
		ChallengeMode:    challengeMode,
		Hypothesis:       hypothesis,
		Type:             lab.Type,
		PathID:           lab.PathID,
		Module:           lab.Module,
		Objectives:       lab.Objectives,
		Lesson:           lesson,
		LockedOnTitles:   lockedOnTitles,
	}
	for _, step := range lab.Loop.Guided {
		view.Guided = append(view.Guided, GuidedStepView{Do: step.Do, Why: step.Why})
	}
	view.Practice.Task = lab.Loop.Practice.Task
	for _, q := range lab.Loop.Practice.Questions {
		view.Practice.Questions = append(view.Practice.Questions, PracticeQuestionView{Prompt: q.Prompt, Answer: q.Answer})
	}
	if haveTicketIdx {
		view.TicketID = webcore.TicketID(ticketIdx)
	}
	if active {
		view.ShellCmd = "docker exec -it breaklab-" + lab.ID + " bash"

		total := len(lab.Hints)
		if v, ok := lab.FindVariant(session.Variant); ok {
			total += len(v.Hints)
		}
		for i := 0; i < total; i++ {
			view.HintSlots = append(view.HintSlots, HintSlotView{Revealed: i < hintsUsed, Cost: hintCost(i)})
		}

		// Diagnostic output only renders for an incident/mystery lab (see
		// lab.html's matching $isMystery gate) — skip the docker call
		// entirely for every other lab type, which never shows it.
		isMystery := lab.Type == "" || lab.Type == "incident"
		if logs != nil && isMystery {
			if out, err := logs(ctx, session.Container, 20); err == nil {
				view.DiagnosticLog = strings.TrimSpace(string(out))
			}
		}
	}
	return view, true, nil
}

// labDetailAPIHandler serves the Astro lab-detail page's data: the
// same BuildDetail computation the old html/template handler used,
// as JSON. There is no Flash here — the API model returns each
// action's own result directly in that action's own response (see
// internal/web/actions) instead of round-tripping it through a
// redirect's query string, so FlashView stays zero-valued.
func labDetailAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cat, err := webcore.LoadCatalog(cfg)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		pf, err := webcore.LoadProgress(cfg)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		session, err := runner.LoadSession(cfg.SessionPath)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		challengeMode := r.URL.Query().Get("mode") == "challenge"

		ctx, cancel := context.WithTimeout(r.Context(), webcore.DockerStatusTimeout)
		defer cancel()

		view, found, err := BuildDetail(ctx, cat, pf, session, r.PathValue("id"), cfg.LabsDir, cfg.ContainerLogs, challengeMode)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !found {
			webcore.WriteJSONError(w, http.StatusNotFound, "lab not found")
			return
		}
		webcore.WriteJSON(w, view)
	}
}

// labLogsStreamHandler streams the active session's own container log
// output as Server-Sent Events, for the lab detail page's live tail.
// It 404s unless the request's {id} matches the active session's lab —
// this route, like every other Docker-backed one, only ever passes
// Docker the container name already recorded in the session file, never
// anything from the request (web-ui spec's "Web UI attack surface
// covers Docker-backed routes" requirement).
func labLogsStreamHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := runner.LoadSession(cfg.SessionPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session == nil || session.Lab != r.PathValue("id") || cfg.ContainerLogsFollow == nil {
			http.NotFound(w, r)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		rc, err := cfg.ContainerLogsFollow(r.Context(), session.Container)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer rc.Close()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		// r.Context() is cancelled when the client disconnects; that's
		// this stream's only stop signal. Closing rc here (rather than
		// only relying on the deferred Close above) is what makes the
		// underlying "docker logs -f" process exit immediately instead
		// of whenever the handler itself next notices — scanner.Scan()
		// below is blocked on a Read that only this Close unblocks.
		go func() {
			<-r.Context().Done()
			_ = rc.Close()
		}()

		// Each line is HTML-escaped and wrapped in a block-level span:
		// htmx's sse-swap inserts this as HTML (not text), and a
		// container's own log output is attacker-influenceable content
		// (anything the broken process prints), so it must never be
		// interpreted as markup in the viewing browser. The span is
		// `display: block` (.log-line in style.css) so successive
		// appended lines render on their own line rather than running
		// together.
		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			line := html.EscapeString(scanner.Text())
			fmt.Fprintf(w, "event: message\ndata: <span class=\"log-line\">%s</span>\n\n", line)
			flusher.Flush()
		}
	}
}

// RegisterRoutes mounts the lab-detail JSON API routes on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/labs/{id}", labDetailAPIHandler(cfg))
	mux.HandleFunc("GET /labs/{id}/logs/stream", labLogsStreamHandler(cfg))
}
