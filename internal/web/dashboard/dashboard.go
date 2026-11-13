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
	"breaklab/internal/web/phases"
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

// StageCardView is one stage's card in the dashboard's stage-card grid
// (add-curriculum-phases): a stage's title and aggregate progress, the
// same per-stage completion buildRoadmap already derives, reused rather
// than recomputed a different way.
type StageCardView struct {
	ID            string
	Title         string
	PhaseID       string
	LabsCompleted int
	LabsTotal     int

	// Complete is true once every lab in this stage has a complete
	// attempt; false for a stage with no labs yet attempted, same as
	// one with some but not all complete — both read as "in progress"
	// for the dashboard's In Progress / Completed tab split.
	Complete bool
}

// DashboardView is the Astro dashboard page's data (GET /api/dashboard).
type DashboardView struct {
	Streak        int
	LabsCompleted int
	VariantsSeen  int
	Active        *ActiveLabView // nil when no lab is running
	Recent        []RecentPostmortemView

	// Stages is one card per regular (non-capstone) stage, for the
	// dashboard's In Progress / Completed stage-card grid.
	Stages []StageCardView

	// Phases is every phase's aggregate progress (including an empty
	// phase, e.g. Phase 2 before it has any stages), for the dashboard's
	// phase-level summary/banner — built by the same function the phase
	// list page uses.
	Phases []phases.PhaseSummaryView

	// MonthlyEngagement counts attempts (started or solved) per day,
	// keyed "2006-01-02", for the current calendar month. A day with no
	// activity is simply absent from the map, not an error — the whole
	// map is empty for a month with no activity at all, rather than nil
	// and omitted.
	MonthlyEngagement map[string]int

	// MostEngagedDay is the "2006-01-02" key with the highest count in
	// MonthlyEngagement, or nil when the month has no activity.
	MostEngagedDay *string

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

// buildStageCards renders one card per regular (non-capstone) stage,
// using the same "does this lab have a complete attempt" definition
// buildDashboard already computes for LabsCompleted.
func buildStageCards(cat *catalog.Catalog, latest map[string]progress.Attempt) []StageCardView {
	regular := cat.RegularStages()
	cards := make([]StageCardView, 0, len(regular))
	for _, group := range regular {
		card := StageCardView{ID: group.Stage.ID, Title: group.Stage.Title, PhaseID: group.Stage.Phase}
		for _, lab := range group.Labs {
			card.LabsTotal++
			if _, ok := latest[lab.ID]; ok {
				card.LabsCompleted++
			}
		}
		card.Complete = card.LabsTotal > 0 && card.LabsCompleted == card.LabsTotal
		cards = append(cards, card)
	}
	return cards
}

// monthlyEngagementDayFormat is the map key format for MonthlyEngagement
// and MostEngagedDay — year-month-day, unambiguous and sortable as a
// string.
const monthlyEngagementDayFormat = "2006-01-02"

// buildMonthlyEngagement counts attempts (by StartedAt, falling back to
// SolvedAt for an attempt with no StartedAt) that fall in now's calendar
// month, one count per day. Returns a non-nil, possibly-empty map (never
// nil) so a month with no activity still renders as "every day empty"
// rather than an omitted field.
func buildMonthlyEngagement(attempts []progress.Attempt, now time.Time) (map[string]int, *string) {
	counts := make(map[string]int)
	year, month, _ := now.Date()

	for _, a := range attempts {
		day := a.StartedAt
		if day.IsZero() && a.SolvedAt != nil {
			day = *a.SolvedAt
		}
		if day.IsZero() {
			continue
		}
		dy, dm, _ := day.Local().Date()
		if dy != year || dm != month {
			continue
		}
		counts[day.Local().Format(monthlyEngagementDayFormat)]++
	}

	days := make([]string, 0, len(counts))
	for day := range counts {
		days = append(days, day)
	}
	sort.Strings(days)

	var mostEngagedDay *string
	best := 0
	for _, day := range days {
		if counts[day] > best {
			best = counts[day]
			d := day
			mostEngagedDay = &d
		}
	}
	return counts, mostEngagedDay
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
	monthlyEngagement, mostEngagedDay := buildMonthlyEngagement(pf.Attempts, time.Now())

	view := DashboardView{
		Streak:              progress.Streak(pf, time.Now()),
		LabsCompleted:       len(latest),
		VariantsSeen:        len(seenVariants),
		Stages:              buildStageCards(cat, latest),
		Phases:              phases.BuildPhaseList(cat, pf).Phases,
		MonthlyEngagement:   monthlyEngagement,
		MostEngagedDay:      mostEngagedDay,
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
