package capstone

import (
	"context"
	"net/http"

	"breaklab/internal/progress"
	"breaklab/internal/runner"
	"breaklab/internal/web/lab"
	"breaklab/internal/web/webcore"
)

// CapstoneView is the Astro capstone page's data (GET /api/capstone). Lab is
// only populated once Unlocked — a locked page shows Progress (how far
// through the regular stages you are) instead of any lab content.
type CapstoneView struct {
	Found     bool // false only if roadmap.yaml declares no capstone stage, or it has no labs yet
	Unlocked  bool
	Progress  []progress.StageStatus // regular-stage completion, shown while locked
	Completed int                    // capstone labs completed, shown once unlocked
	Total     int
	Lab       lab.LabDetailView
}

// buildCapstone assembles CapstoneView for the single capstone lab (the
// first under the catalog's capstone stage, if any).
func buildCapstone(ctx context.Context, cfg webcore.Config, pf *progress.ProgressFile, session *runner.Session) (CapstoneView, error) {
	cat, err := webcore.LoadCatalog(cfg)
	if err != nil {
		return CapstoneView{}, err
	}

	group := cat.CapstoneStage()
	if group == nil || len(group.Labs) == 0 {
		return CapstoneView{}, nil
	}

	status := progress.Capstone(cat, pf)
	view := CapstoneView{
		Found:     true,
		Unlocked:  status != nil && status.Unlocked,
		Progress:  progress.StageCompletion(cat, pf),
		Completed: status.Completed,
		Total:     status.Total,
	}
	if !view.Unlocked {
		return view, nil
	}

	detail, found, err := lab.BuildDetail(ctx, cat, pf, session, group.Labs[0].ID, cfg.LabsDir, cfg.ContainerLogs, false)
	if err != nil {
		return CapstoneView{}, err
	}
	view.Found = found
	view.Lab = detail
	return view, nil
}

// capstoneAPIHandler serves the Astro capstone page's data: the same
// buildCapstone computation the old html/template handler used, as
// JSON. Returns 404 when !Found (no capstone stage declared, or it has
// no labs yet), the same as the old handler.
func capstoneAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		view, err := buildCapstone(ctx, cfg, pf, session)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !view.Found {
			webcore.WriteJSONError(w, http.StatusNotFound, "no capstone stage")
			return
		}
		webcore.WriteJSON(w, view)
	}
}

// RegisterRoutes mounts the capstone JSON API route on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/capstone", capstoneAPIHandler(cfg))
}
