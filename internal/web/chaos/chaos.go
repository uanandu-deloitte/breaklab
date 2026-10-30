package chaos

import (
	"net/http"
	"sort"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
	"breaklab/internal/web/webcore"
)

// chaosConcepts groups existing labs by chaos-engineering principle,
// cross-cutting roadmap.yaml's stage order rather than replacing it — a
// lab here still belongs to its one real stage, it's just also matched
// by tag into one of these concepts.
var chaosConcepts = []struct {
	ID    string
	Title string
	Blurb string
	Tag   string
}{
	{"fault-injection", "Fault injection", "Kill or starve a process and watch what actually happens.", "fault-injection"},
	{"dependency-failure", "Dependency failure", "The thing you depend on dies; your own process never notices.", "dependency-failure"},
	{"blind-spots", "Blind spots", "Chaos isn't just outages — it's monitoring that quietly lies to you.", "blind-spots"},
	{"blast-radius", "Blast radius", "Contain a failure before it becomes systemic.", "blast-radius"},
}

// ChaosConceptView is one chaos-engineering principle and the labs
// tagged under it, in difficulty order.
type ChaosConceptView struct {
	ID    string
	Title string
	Blurb string
	Labs  []webcore.LabRowView
}

// ChaosView is the Astro chaos page's data (GET /api/chaos).
type ChaosView struct {
	Concepts []ChaosConceptView
}

// buildChaos groups every chaos-tagged lab, across every stage, by
// chaosConcepts' tag — real progress status, no lock/up-next semantics
// since this is a browsable reference, not a gated path.
func buildChaos(cat *catalog.Catalog, pf *progress.ProgressFile) ChaosView {
	progress.ApplyStatus(cat, pf)
	latest := progress.LatestCompleteAttempts(pf)
	reviewDue := progress.ReviewDue(pf)

	view := ChaosView{Concepts: make([]ChaosConceptView, 0, len(chaosConcepts))}
	for _, concept := range chaosConcepts {
		cv := ChaosConceptView{ID: concept.ID, Title: concept.Title, Blurb: concept.Blurb}
		for _, group := range cat.RegularStages() {
			for _, lab := range group.Labs {
				if webcore.HasTag(lab.Tags, concept.Tag) {
					cv.Labs = append(cv.Labs, webcore.LabRow(group, lab, latest, pf, reviewDue))
				}
			}
		}
		sort.Slice(cv.Labs, func(i, j int) bool {
			if cv.Labs[i].Difficulty != cv.Labs[j].Difficulty {
				return cv.Labs[i].Difficulty < cv.Labs[j].Difficulty
			}
			return cv.Labs[i].ID < cv.Labs[j].ID
		})
		view.Concepts = append(view.Concepts, cv)
	}
	return view
}

// chaosAPIHandler serves the Astro chaos page's data: the same
// buildChaos computation the old html/template handler used, as JSON.
func chaosAPIHandler(cfg webcore.Config) http.HandlerFunc {
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

		view := buildChaos(cat, pf)
		webcore.WriteJSON(w, view)
	}
}

// RegisterRoutes mounts the chaos JSON API route on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/chaos", chaosAPIHandler(cfg))
}
