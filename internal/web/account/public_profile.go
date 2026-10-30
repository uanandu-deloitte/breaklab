package account

import (
	"net/http"

	"breaklab/internal/progress"
	"breaklab/internal/site"
	"breaklab/internal/web/webcore"
)

// publicProfileView is GET /api/u/{name}'s data — site.IndexView (the
// same shape breaklab publish's static site renders) plus the
// account's username for the page title.
type publicProfileView struct {
	Username string
	Stages   []site.StageView
	Totals   site.Totals
}

// publicProfileAPIHandler serves the opt-in /u/{name} page's data from
// name's hosted (synced) attempts. No postmortem content is ever
// synced to the hosted side (only Attempt records, which merely name a
// local file path) — passing an empty "published" map to
// site.BuildIndex means no postmortem is ever linked or named here,
// satisfying "nothing leaks" by construction rather than by checking a
// flag this server has no way to see.
func publicProfileAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		u, err := cfg.AccountStore.UserByUsername(name)
		if err != nil || !u.PublicProfile {
			webcore.WriteJSONError(w, http.StatusNotFound, "not found")
			return
		}

		cat, err := webcore.LoadCatalog(cfg)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		hosted, err := cfg.AccountStore.PullAttempts(u.ID)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		index := site.BuildIndex(cat, &progress.ProgressFile{Attempts: hosted}, map[string]bool{})
		webcore.WriteJSON(w, publicProfileView{
			Username: u.Username,
			Stages:   index.Stages,
			Totals:   index.Totals,
		})
	}
}
