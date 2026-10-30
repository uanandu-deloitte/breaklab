package challenges

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
)

// ChallengeFilters echoes the current filter values back to the form.
type ChallengeFilters struct {
	Stage      string
	Difficulty string
	Status     string
	Tag        string
	Query      string // free-text search over title and tags
	Sort       string // "recent" (default), "difficulty", or "alpha"
}

// ActiveIncidentView is the challenge browser's "a lab is active"
// banner.
type ActiveIncidentView struct {
	Lab      string
	TicketID string
	Href     string
}

// ChallengeView is the Astro challenge-browser page's data (GET /api/challenges).
type ChallengeView struct {
	Rows           []webcore.LabRowView
	Filters        ChallengeFilters
	Stages         []catalog.Stage // for the stage filter's options, in roadmap order
	Tags           []string        // for the topic filter's options, sorted
	ActiveIncident *ActiveIncidentView
	Velocity       progress.TrainingVelocity

	// SandboxImagesSize is Docker's real image-store disk usage,
	// formatted for display; empty when Docker is unreachable.
	SandboxImagesSize string
}

// formatBytes renders n as a short "1.20 GiB"-style label for the
// challenge browser's sandbox-image-cache tile.
func formatBytes(n uint64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := uint64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}

// catalogTags collects every distinct tag used across the catalog, sorted,
// for the topic filter's dropdown options.
func catalogTags(cat *catalog.Catalog) []string {
	seen := make(map[string]bool)
	for _, group := range cat.RegularStages() {
		for _, lab := range group.Labs {
			for _, tag := range lab.Tags {
				seen[tag] = true
			}
		}
	}
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// matchedLab is one lab that passed buildChallenges'/the random-launch
// route's filters, paired with its stage group (webcore.LabRow needs both).
type matchedLab struct {
	group catalog.StageGroup
	lab   catalog.Lab
}

// matchesQuery reports whether lab's title or any tag contains query
// (already lowercased), for the challenge browser's free-text search.
func matchesQuery(lab catalog.Lab, query string) bool {
	if strings.Contains(strings.ToLower(lab.Title), query) {
		return true
	}
	for _, tag := range lab.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

// filterLabs returns every lab matching filters, across every stage.
// Shared by buildChallenges (which also sorts and renders the result)
// and the random-launch route (which just needs the matching set to
// pick from) — both apply exactly the same rules, so a lab the browser
// is showing can always be the one random-launch picks.
func filterLabs(cat *catalog.Catalog, filters ChallengeFilters) []matchedLab {
	var difficulty int
	difficultyFilterActive := filters.Difficulty != ""
	if difficultyFilterActive {
		difficulty, _ = strconv.Atoi(filters.Difficulty) // 0 on error: matches nothing, no error surfaced
	}
	query := strings.ToLower(strings.TrimSpace(filters.Query))

	var matched []matchedLab
	for _, group := range cat.RegularStages() {
		for _, lab := range group.Labs {
			if filters.Stage != "" && group.Stage.ID != filters.Stage {
				continue
			}
			if difficultyFilterActive && lab.Difficulty != difficulty {
				continue
			}
			if filters.Status != "" && lab.Status != filters.Status {
				continue
			}
			if filters.Tag != "" && !webcore.HasTag(lab.Tags, filters.Tag) {
				continue
			}
			if query != "" && !matchesQuery(lab, query) {
				continue
			}
			matched = append(matched, matchedLab{group, lab})
		}
	}
	return matched
}

// lastActivity maps each lab id to its most recently started attempt's
// StartedAt (any status), for the "recent activity" sort order. A lab
// with no recorded attempt is the zero time, which sorts last.
func lastActivity(pf *progress.ProgressFile) map[string]time.Time {
	activity := make(map[string]time.Time)
	for _, a := range pf.Attempts {
		if a.StartedAt.After(activity[a.Lab]) {
			activity[a.Lab] = a.StartedAt
		}
	}
	return activity
}

// sortMatched orders matched in place according to sortBy ("difficulty",
// "alpha", or anything else — including "" — for "recent activity",
// the default).
func sortMatched(matched []matchedLab, activity map[string]time.Time, sortBy string) {
	switch sortBy {
	case "difficulty":
		sort.SliceStable(matched, func(i, j int) bool {
			return matched[i].lab.Difficulty < matched[j].lab.Difficulty
		})
	case "alpha":
		sort.SliceStable(matched, func(i, j int) bool {
			return strings.ToLower(matched[i].lab.Title) < strings.ToLower(matched[j].lab.Title)
		})
	default:
		sort.SliceStable(matched, func(i, j int) bool {
			return activity[matched[i].lab.ID].After(activity[matched[j].lab.ID])
		})
	}
}

// buildChallenges lists every lab across every stage, narrowed by the
// stage/difficulty/status/tag/search query params and ordered by sort.
// An unrecognized or unparseable filter value yields an empty result
// (no lab can match it) rather than an error; no params (or clearing
// them) returns the full list. diskUsage is called only to populate the
// sandbox-image-cache tile; a nil func or a lookup error just leaves it
// blank rather than failing the page (same degrade-gracefully contract
// as topology's status func).
func buildChallenges(ctx context.Context, cat *catalog.Catalog, pf *progress.ProgressFile, session *runner.Session, q url.Values, diskUsage webcore.DiskUsageFunc) ChallengeView {
	progress.ApplyStatus(cat, pf)
	latest := progress.LatestCompleteAttempts(pf)
	reviewDue := progress.ReviewDue(pf)

	filters := ChallengeFilters{
		Stage:      q.Get("stage"),
		Difficulty: q.Get("difficulty"),
		Status:     q.Get("status"),
		Tag:        q.Get("tag"),
		Query:      q.Get("q"),
		Sort:       q.Get("sort"),
	}

	matched := filterLabs(cat, filters)
	sortMatched(matched, lastActivity(pf), filters.Sort)

	rows := make([]webcore.LabRowView, 0, len(matched))
	for _, m := range matched {
		rows = append(rows, webcore.LabRow(m.group, m.lab, latest, pf, reviewDue))
	}

	regular := cat.RegularStages()
	stages := make([]catalog.Stage, 0, len(regular))
	for _, group := range regular {
		stages = append(stages, group.Stage)
	}

	view := ChallengeView{
		Rows:     rows,
		Filters:  filters,
		Stages:   stages,
		Tags:     catalogTags(cat),
		Velocity: progress.BuildTrainingVelocity(cat, pf),
	}

	if session != nil {
		view.ActiveIncident = &ActiveIncidentView{
			Lab:      session.Lab,
			TicketID: webcore.TicketID(session.AttemptIndex),
			Href:     "/labs/" + session.Lab,
		}
	}

	if diskUsage != nil {
		if du, err := diskUsage(ctx); err == nil {
			view.SandboxImagesSize = formatBytes(du.ImagesSizeBytes)
		}
	}

	return view
}

// challengesAPIHandler serves the Astro challenge-browser page's data:
// the same buildChallenges computation the old html/template handler
// used, as JSON.
func challengesAPIHandler(cfg webcore.Config) http.HandlerFunc {
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

		ctx, cancel := context.WithTimeout(r.Context(), webcore.DockerStatusTimeout)
		defer cancel()

		view := buildChallenges(ctx, cat, pf, session, r.URL.Query(), cfg.DiskUsage)
		webcore.WriteJSON(w, view)
	}
}

// randomChallengePayload is challengesRandomAPIHandler's response: the
// randomly-picked lab's id, or "" when nothing matches the filters.
type randomChallengePayload struct {
	Lab string `json:"lab"`
}

// challengesRandomAPIHandler is the JSON equivalent of
// challengesRandomHandler: same filterLabs rules and random pick,
// returned as {"lab": "<id>"} instead of a redirect. It never starts a
// lab itself — there remains exactly one place a lab actually starts,
// the detail page's own Launch action.
func challengesRandomAPIHandler(cfg webcore.Config) http.HandlerFunc {
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
		progress.ApplyStatus(cat, pf)

		q := r.URL.Query()
		filters := ChallengeFilters{
			Stage:      q.Get("stage"),
			Difficulty: q.Get("difficulty"),
			Status:     q.Get("status"),
			Tag:        q.Get("tag"),
			Query:      q.Get("q"),
		}
		matched := filterLabs(cat, filters)
		if len(matched) == 0 {
			webcore.WriteJSON(w, randomChallengePayload{})
			return
		}
		pick := matched[rand.IntN(len(matched))]
		webcore.WriteJSON(w, randomChallengePayload{Lab: pick.lab.ID})
	}
}

// RegisterRoutes mounts the challenge-browser JSON API routes on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/challenges", challengesAPIHandler(cfg))
	mux.HandleFunc("POST /api/challenges/random", challengesRandomAPIHandler(cfg))
}
