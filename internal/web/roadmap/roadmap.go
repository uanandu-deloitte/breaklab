package roadmap

import (
	"net/http"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/web/webcore"
)

// RoadmapStageView is one roadmap stage and its labs, in lab.yaml order.
type RoadmapStageView struct {
	ID    string
	Title string
	Blurb string // this stage's own "why it's next" line, shown before it
	Index int    // 1-based position in roadmap.yaml order, for the "stage N" ribbon
	Labs  []webcore.LabRowView

	// Track is catalog.Stage.Track, carried through unchanged: non-empty
	// only when this stage is one of a set of alternatives a learner
	// picks between (see catalog.GroupByTrack). Empty for every ordinary
	// stage.
	Track string

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

// BuildStageViews renders groups (already stage-status-applied via
// progress.ApplyStatus) into RoadmapStageViews: 1-based Index, connector
// blurbs, and the single "up next" cursor are all scoped to just this
// slice, so a caller can reuse it for either the full roadmap or one
// phase's stages. Exported for reuse by internal/web/phases and
// internal/web/dashboard.
func BuildStageViews(groups []catalog.StageGroup, pf *progress.ProgressFile, latest map[string]progress.Attempt, reviewDue map[string]bool) []RoadmapStageView {
	stages := make([]RoadmapStageView, 0, len(groups))
	for i, group := range groups {
		sv := RoadmapStageView{
			ID:    group.Stage.ID,
			Title: group.Stage.Title,
			Blurb: group.Stage.Blurb,
			Index: i + 1,
			Track: group.Stage.Track,
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
	return stages
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

// RegisterRoutes mounts the stage-skip JSON API route on mux, under the
// /api/phases namespace — the only UI surface left that calls it is
// PhaseDetail.tsx.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("POST /api/phases/stages/{id}/skip", stageSkipAPIHandler(cfg))
}
