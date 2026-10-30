package roadmap

import (
	"net/http"
	"runtime/debug"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/web/webcore"
)

// buildRevision returns the short VCS revision Go embeds automatically
// when the binary was built from a git checkout (`go build`, not `go
// run`), for the roadmap footer's real build-identity line. Empty when
// unavailable — never a fabricated placeholder.
func buildRevision() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	var revision string
	var modified bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if revision == "" {
		return ""
	}
	if len(revision) > 7 {
		revision = revision[:7]
	}
	if modified {
		revision += "-dirty"
	}
	return revision
}

// RoadmapStageView is one roadmap stage and its labs, in lab.yaml order.
type RoadmapStageView struct {
	ID    string
	Title string
	Blurb string // this stage's own "why it's next" line, shown before it
	Index int    // 1-based position in roadmap.yaml order, for the "stage N" ribbon
	Labs  []webcore.LabRowView

	// ShowSkip is true while at least one lab in this stage isn't yet
	// LabStatusComplete — the condition for showing the "I already know
	// this" skip action. A fully complete stage has nothing left to skip.
	ShowSkip bool

	// ConnectorBlurb is the *next* stage's Blurb, shown beside the arrow
	// leading into it — set on every stage but the last, which has no
	// connector after it.
	ConnectorBlurb string
	IsLast         bool
}

// RoadmapView is the Astro roadmap page's data (GET /api/roadmap).
type RoadmapView struct {
	Stages        []RoadmapStageView
	BuildRevision string // real VCS revision Go embeds at build time; "" when unavailable (e.g. `go run`)
}

// buildRoadmap renders every stage from roadmap.yaml in order, each with
// its labs in lab.yaml order.
func buildRoadmap(cat *catalog.Catalog, pf *progress.ProgressFile) RoadmapView {
	progress.ApplyStatus(cat, pf)
	latest := progress.LatestCompleteAttempts(pf)
	reviewDue := progress.ReviewDue(pf)

	regular := cat.RegularStages()
	stages := make([]RoadmapStageView, 0, len(regular))
	for i, group := range regular {
		sv := RoadmapStageView{
			ID:    group.Stage.ID,
			Title: group.Stage.Title,
			Blurb: group.Stage.Blurb,
			Index: i + 1,
		}
		for _, lab := range group.Labs {
			row := webcore.LabRow(group, lab, latest, pf, reviewDue)
			if row.Status != progress.LabStatusComplete {
				sv.ShowSkip = true
			}
			sv.Labs = append(sv.Labs, row)
		}
		stages = append(stages, sv)
	}
	for i := range stages {
		if i == len(stages)-1 {
			stages[i].IsLast = true
			continue
		}
		stages[i].ConnectorBlurb = stages[i+1].Blurb
	}
	markUpNext(stages)
	return RoadmapView{Stages: stages, BuildRevision: buildRevision()}
}

// markUpNext sets IsUpNext on the single lab, across every stage in
// order, that is the first one not yet "complete" — the room-list's one
// "Up next" cursor. A roadmap where every lab is complete marks none.
func markUpNext(stages []RoadmapStageView) {
	for i := range stages {
		for j := range stages[i].Labs {
			status := stages[i].Labs[j].Status
			// A locked lab can't be started yet, so it's never the "Up
			// next" cursor — skip past it to the first lab that's
			// actually complete-or-startable.
			if status == progress.LabStatusComplete || status == progress.LabStatusLocked {
				continue
			}
			stages[i].Labs[j].IsUpNext = true
			return
		}
	}
}

// stageSkipAPIHandler marks every not-yet-attempted lab in the path
// stage id as skipped (progress.SkipStage) — the "I already know
// this" action next to a stage's heading, with a JSON response
// instead of a redirect.
func stageSkipAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		cat, err := webcore.LoadCatalog(cfg)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var labIDs []string
		for _, group := range cat.RegularStages() {
			if group.Stage.ID != id {
				continue
			}
			for _, lab := range group.Labs {
				labIDs = append(labIDs, lab.ID)
			}
			break
		}

		if err := progress.SkipStage(cfg.ProgressPath, labIDs); err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// roadmapAPIHandler serves the Astro roadmap page's data: the same
// buildRoadmap computation and webcore.ProgressForRequest session-aware
// merge the old html/template handler used, as JSON.
func roadmapAPIHandler(cfg webcore.Config) http.HandlerFunc {
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

		view := buildRoadmap(cat, pf)
		webcore.WriteJSON(w, view)
	}
}

// RegisterRoutes mounts the roadmap JSON API routes on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/roadmap", roadmapAPIHandler(cfg))
	mux.HandleFunc("POST /api/roadmap/stages/{id}/skip", stageSkipAPIHandler(cfg))
}
