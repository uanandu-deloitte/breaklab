// Package phases serves the Phase browsing views (GET /api/phases, GET
// /api/phases/{id}): the grouping roadmap.yaml declares above Stage, so
// the web UI can show one phase at a time instead of a single flat
// all-stages list.
package phases

import (
	"net/http"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/web/roadmap"
	"breaklab/internal/web/webcore"
)

// PhaseSummaryView is one phase's entry in the phase list: aggregate
// progress across its stages, not its stages themselves.
type PhaseSummaryView struct {
	ID            string
	Title         string
	Blurb         string
	LabsCompleted int
	LabsTotal     int
}

// PhaseListView is the Astro phase-list page's data (GET /api/phases).
type PhaseListView struct {
	Phases []PhaseSummaryView
}

// PhaseDetailView is one phase's detail page data (GET /api/phases/{id}):
// its stages, each with its labs, in roadmap order. A phase with no
// stages yet still renders, with an empty Stages slice.
type PhaseDetailView struct {
	ID     string
	Title  string
	Blurb  string
	Stages []roadmap.RoadmapStageView
}

// summarize computes one phase's aggregate labs-completed/total from its
// (status-applied) stage groups, counting only the regular (non-capstone)
// stages — the same scope the phase detail page's Stages uses.
func summarize(group catalog.PhaseGroup) PhaseSummaryView {
	sv := PhaseSummaryView{ID: group.Phase.ID, Title: group.Phase.Title, Blurb: group.Phase.Blurb}
	for _, stage := range group.RegularStages() {
		for _, lab := range stage.Labs {
			sv.LabsTotal++
			if lab.Status == progress.LabStatusComplete {
				sv.LabsCompleted++
			}
		}
	}
	return sv
}

// BuildPhaseList renders every phase from roadmap.yaml in order with its
// aggregate progress. Exported for reuse by internal/web/dashboard, so
// the dashboard's phase summaries and the phase list page never compute
// "phase progress" two different ways.
func BuildPhaseList(cat *catalog.Catalog, pf *progress.ProgressFile) PhaseListView {
	progress.ApplyStatus(cat, pf)

	view := PhaseListView{Phases: make([]PhaseSummaryView, 0, len(cat.Phases))}
	for _, group := range cat.Phases {
		view.Phases = append(view.Phases, summarize(group))
	}
	return view
}

// buildPhaseDetail renders one phase's stages (regular, non-capstone —
// matching the existing roadmap view's scope) in roadmap order, each
// with its labs. found is false when phaseID names no declared phase.
func buildPhaseDetail(cat *catalog.Catalog, pf *progress.ProgressFile, phaseID string) (PhaseDetailView, bool) {
	progress.ApplyStatus(cat, pf)

	group, ok := cat.FindPhase(phaseID)
	if !ok {
		return PhaseDetailView{}, false
	}

	latest := progress.LatestCompleteAttempts(pf)
	reviewDue := progress.ReviewDue(pf)
	stages := roadmap.BuildStageViews(group.RegularStages(), pf, latest, reviewDue)

	return PhaseDetailView{
		ID:     group.Phase.ID,
		Title:  group.Phase.Title,
		Blurb:  group.Phase.Blurb,
		Stages: stages,
	}, true
}

// phaseListAPIHandler serves the Astro phase-list page's data.
func phaseListAPIHandler(cfg webcore.Config) http.HandlerFunc {
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

		webcore.WriteJSON(w, BuildPhaseList(cat, pf))
	}
}

// phaseDetailAPIHandler serves one phase's detail page data.
func phaseDetailAPIHandler(cfg webcore.Config) http.HandlerFunc {
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

		view, found := buildPhaseDetail(cat, pf, r.PathValue("id"))
		if !found {
			webcore.WriteJSONError(w, http.StatusNotFound, "phase not found")
			return
		}
		webcore.WriteJSON(w, view)
	}
}

// RegisterRoutes mounts the phase JSON API routes on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/phases", phaseListAPIHandler(cfg))
	mux.HandleFunc("GET /api/phases/{id}", phaseDetailAPIHandler(cfg))
}
