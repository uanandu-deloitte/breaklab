package dashboard

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
)

const recentPostmortemsLimit = 5

// ActiveLabView is the dashboard's "a lab is running" block.
type ActiveLabView struct {
	Lab            string
	Title          string
	Href           string
	TicketID       string
	ElapsedMinutes int
	HintsUsed      int
	HintsTotal     int
	Ports          []webcore.PortNodeView
	StageTitle     string
	Tags           []string
	ShellCmd       string
	DiagnosticLog  string // tail of the container's own logs, when available
	ReminderCmd    string
}

// RecentPostmortemView is one entry in the dashboard's recent-postmortems list.
type RecentPostmortemView struct {
	Lab      string
	TicketID string
	Date     string
	Href     string
}

// DashboardView is the Astro dashboard page's data (GET /api/dashboard).
type DashboardView struct {
	Streak        int
	LabsCompleted int
	VariantsSeen  int
	Active        *ActiveLabView // nil when no lab is running
	Recent        []RecentPostmortemView

	// DockerVersion and Uptime are real telemetry, degrading to "" when
	// Docker is slow/unreachable or the server just started — never a
	// fabricated value.
	DockerVersion       string
	Uptime              string
	AvgResolutionTime   string
	PostmortemDrafts    int
	PostmortemPublished int

	// Patterns lists topic tags needing consistently more hints than
	// average (progress.Patterns); empty when none qualify.
	Patterns []progress.Pattern
}

// buildDashboard assembles DashboardView. LabsCompleted/VariantsSeen use
// the same complete-only, any-status definitions as the public site
// (design.md Decision 6). Recent is the five most recently solved
// complete attempts, regardless of their postmortem's publish flag
// (design.md's "Spec clarifications," item 2). version and logs are
// called only with session.Container (never request input, there being
// no request here); a nil func or a lookup error just leaves the
// corresponding field at its zero value rather than failing the page,
// the same degrade-gracefully contract as topology's status func.
func buildDashboard(ctx context.Context, cat *catalog.Catalog, pf *progress.ProgressFile, session *runner.Session, started time.Time, version webcore.VersionFunc, logs webcore.LogsFunc) DashboardView {
	latest := progress.LatestCompleteAttempts(pf)

	seenVariants := make(map[string]bool)
	for _, a := range pf.Attempts {
		seenVariants[a.Lab+"\x00"+a.Variant] = true
	}

	drafts, published := progress.PostmortemCounts(pf.Attempts)

	view := DashboardView{
		Streak:              progress.Streak(pf, time.Now()),
		LabsCompleted:       len(latest),
		VariantsSeen:        len(seenVariants),
		AvgResolutionTime:   progress.FormatDuration(progress.AvgResolutionTime(pf.Attempts)),
		PostmortemDrafts:    drafts,
		PostmortemPublished: published,
		Patterns:            progress.Patterns(cat, pf),
	}

	if version != nil {
		if v, err := version(ctx); err == nil {
			view.DockerVersion = v
		}
	}
	if !started.IsZero() {
		view.Uptime = progress.FormatDuration(time.Since(started))
	}

	if session != nil {
		var elapsed int
		if session.AttemptIndex >= 0 && session.AttemptIndex < len(pf.Attempts) {
			elapsed = int(time.Since(pf.Attempts[session.AttemptIndex].StartedAt).Minutes())
		}
		active := &ActiveLabView{
			Lab:            session.Lab,
			Href:           "/labs/" + session.Lab,
			TicketID:       webcore.TicketID(session.AttemptIndex),
			ElapsedMinutes: elapsed,
			HintsUsed:      session.HintsUsed,
			ReminderCmd:    "breaklab check",
			ShellCmd:       "docker exec -it " + session.Container + " bash",
		}
		if lab, stageTitle, found := cat.FindLab(session.Lab); found {
			active.Title = lab.Title
			active.StageTitle = stageTitle
			active.Tags = lab.Tags
			active.HintsTotal = len(lab.Hints)
			if v, ok := lab.FindVariant(session.Variant); ok {
				active.HintsTotal += len(v.Hints)
			}
			for _, p := range lab.Environment.Ports {
				host, container, ok := strings.Cut(p, ":")
				pv := webcore.PortNodeView{HostPort: host, ContainerPort: container}
				if ok && host != "" {
					pv.OpenHref = "http://localhost:" + host
				}
				active.Ports = append(active.Ports, pv)
			}
		}
		if logs != nil {
			if out, err := logs(ctx, session.Container, 20); err == nil {
				active.DiagnosticLog = strings.TrimSpace(string(out))
			}
		}
		view.Active = active
	}

	type labAttempt struct {
		lab     string
		attempt progress.Attempt
	}
	entries := make([]labAttempt, 0, len(latest))
	for lab, a := range latest {
		if a.Postmortem == "" { // never written, or rejected by webcore.LoadProgress
			continue
		}
		entries = append(entries, labAttempt{lab, a})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].attempt.SolvedAt.After(*entries[j].attempt.SolvedAt)
	})
	if len(entries) > recentPostmortemsLimit {
		entries = entries[:recentPostmortemsLimit]
	}
	for _, e := range entries {
		recent := RecentPostmortemView{
			Lab:  e.lab,
			Date: e.attempt.SolvedAt.Local().Format("2006-01-02"),
			Href: "/postmortems/" + webcore.SlugFor(e.attempt.Postmortem),
		}
		if idx, ok := attemptIndexOf(pf, e.attempt); ok {
			recent.TicketID = webcore.TicketID(idx)
		}
		view.Recent = append(view.Recent, recent)
	}

	return view
}

// attemptIndexOf finds target's position within pf.Attempts, matching
// on lab+variant+start time (unique per attempt in practice). Used to
// recover the index a progress.LatestCompleteAttempts copy came from,
// since TicketID needs the index, not the attempt value.
func attemptIndexOf(pf *progress.ProgressFile, target progress.Attempt) (int, bool) {
	for i, a := range pf.Attempts {
		if a.Lab == target.Lab && a.Variant == target.Variant && a.StartedAt.Equal(target.StartedAt) {
			return i, true
		}
	}
	return -1, false
}

// dashboardAPIHandler serves the Astro dashboard page's data: the same
// buildDashboard computation and webcore.ProgressForRequest session-aware
// merge the old html/template handler used, as JSON.
func dashboardAPIHandler(cfg webcore.Config) http.HandlerFunc {
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
		pf = webcore.ProgressForRequest(cfg, r, pf)
		session, err := runner.LoadSession(cfg.SessionPath)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), webcore.DockerStatusTimeout)
		defer cancel()

		view := buildDashboard(ctx, cat, pf, session, cfg.Started, cfg.DockerVersion, cfg.ContainerLogs)
		webcore.WriteJSON(w, view)
	}
}

// RegisterRoutes mounts the dashboard JSON API route on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/dashboard", dashboardAPIHandler(cfg))
}
