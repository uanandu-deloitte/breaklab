package postmortems

import (
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"breaklab/internal/catalog"
	"breaklab/internal/postmortem"
	"breaklab/internal/progress"
	"breaklab/internal/web/webcore"
)

// SectionView is one rendered section of a postmortem.
type SectionView struct {
	Name   string
	Blocks []postmortem.Block
}

// PostmortemPageView is the Astro postmortem-detail page's data
// (GET /api/postmortems/{slug}). Unlike the public site, this always
// renders regardless of the file's publish flag — that flag only ever
// controlled what becomes public (design.md's Postmortem viewer
// requirement).
type PostmortemPageView struct {
	Lab      string
	Variant  string
	Date     string
	TicketID string
	Sections []SectionView
}

// postmortemSlugs maps each known postmortem's URL slug to its actual
// file path, built only from "complete" attempts' Postmortem fields.
// A request slug that isn't a key here is a 404 — never a file read —
// so there is no path-traversal surface by construction (design.md
// Decision 7).
func postmortemSlugs(pf *progress.ProgressFile) map[string]string {
	slugs := make(map[string]string)
	for _, a := range pf.Attempts {
		if a.Status == progress.StatusComplete && a.Postmortem != "" {
			slugs[webcore.SlugFor(a.Postmortem)] = a.Postmortem
		}
	}
	return slugs
}

// attemptIndexForPostmortemPath finds the index within pf.Attempts of
// the attempt whose (already-confined, by webcore.LoadProgress) Postmortem
// path equals path, for the postmortem page's ticket-id header.
func attemptIndexForPostmortemPath(pf *progress.ProgressFile, path string) (int, bool) {
	for i, a := range pf.Attempts {
		if a.Postmortem == path {
			return i, true
		}
	}
	return -1, false
}

func buildPostmortemPage(path string, cat *catalog.Catalog) (PostmortemPageView, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PostmortemPageView{}, err
	}
	raw, body, err := postmortem.SplitFrontMatter(string(data))
	if err != nil {
		return PostmortemPageView{}, err
	}
	fm, err := postmortem.ParseFrontMatter(raw)
	if err != nil {
		return PostmortemPageView{}, err
	}

	labName := fm.Lab
	if lab, _, ok := cat.FindLab(fm.Lab); ok {
		labName = lab.Title
	}

	rawSections := postmortem.SplitIntoSections(body)
	sections := make([]SectionView, 0, len(postmortem.DisplaySections))
	for _, name := range postmortem.DisplaySections {
		text := postmortem.StripComments(rawSections[name])
		// Hypothesis is optional (unlike the required sections it's
		// listed alongside): an older postmortem with none, or one
		// where the learner left it blank, just omits the section
		// instead of showing an empty "## Hypothesis" heading.
		if name == "Hypothesis" && strings.TrimSpace(text) == "" {
			continue
		}
		sections = append(sections, SectionView{Name: name, Blocks: postmortem.SplitBlocks(text)})
	}

	return PostmortemPageView{Lab: labName, Variant: fm.Variant, Date: fm.Date, Sections: sections}, nil
}

// isPublished reports whether the postmortem at path has publish: true
// in its front matter. A missing, unreadable, or malformed file is
// treated as not published, the same degrade-to-draft rule
// internal/site's LoadPublished applies for the public site.
func isPublished(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	raw, _, err := postmortem.SplitFrontMatter(string(data))
	if err != nil {
		return false
	}
	fm, err := postmortem.ParseFrontMatter(raw)
	if err != nil {
		return false
	}
	return fm.Publish
}

// PostmortemRowView is one row in the postmortems index.
type PostmortemRowView struct {
	Lab        string
	StageTitle string
	Difficulty int
	Date       string
	Status     string // "draft" or "published"
	TicketID   string
	Href       string
}

// PostmortemsIndexFilters echoes the current filter values back to the form.
type PostmortemsIndexFilters struct {
	Stage      string
	Difficulty string
	Status     string
}

// PostmortemsIndexView is the Astro postmortems-index page's data
// (GET /api/postmortems).
type PostmortemsIndexView struct {
	Rows    []PostmortemRowView
	Filters PostmortemsIndexFilters
	Stages  []catalog.Stage // for the stage filter's options, in roadmap order
}

// buildPostmortemsIndex lists every postmortem the user has written —
// draft and published alike — narrowed by the stage/difficulty/status
// query params, most recently solved first. Only the latest complete
// attempt per lab is considered, the same rule the challenge browser and
// BuildIndex already apply.
func buildPostmortemsIndex(cat *catalog.Catalog, pf *progress.ProgressFile, q url.Values) PostmortemsIndexView {
	latest := progress.LatestCompleteAttempts(pf)

	filters := PostmortemsIndexFilters{
		Stage:      q.Get("stage"),
		Difficulty: q.Get("difficulty"),
		Status:     q.Get("status"),
	}
	var difficulty int
	difficultyActive := filters.Difficulty != ""
	if difficultyActive {
		difficulty, _ = strconv.Atoi(filters.Difficulty) // 0 on error: matches nothing, no error surfaced
	}

	regular := cat.RegularStages()
	stages := regular
	if capstone := cat.CapstoneStage(); capstone != nil {
		stages = append(append([]catalog.StageGroup{}, regular...), *capstone)
	}

	var rows []PostmortemRowView
	for _, group := range stages {
		for _, lab := range group.Labs {
			attempt, ok := latest[lab.ID]
			if !ok || attempt.Postmortem == "" {
				continue
			}
			if filters.Stage != "" && group.Stage.ID != filters.Stage {
				continue
			}
			if difficultyActive && lab.Difficulty != difficulty {
				continue
			}
			status := "draft"
			if isPublished(attempt.Postmortem) {
				status = "published"
			}
			if filters.Status != "" && filters.Status != status {
				continue
			}

			var ticket string
			if idx, ok := attemptIndexForPostmortemPath(pf, attempt.Postmortem); ok {
				ticket = webcore.TicketID(idx)
			}
			rows = append(rows, PostmortemRowView{
				Lab:        lab.Title,
				StageTitle: group.Stage.Title,
				Difficulty: lab.Difficulty,
				Date:       attempt.SolvedAt.Local().Format("2006-01-02"),
				Status:     status,
				TicketID:   ticket,
				Href:       "/postmortems/" + webcore.SlugFor(attempt.Postmortem),
			})
		}
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Date > rows[j].Date })

	stageOptions := make([]catalog.Stage, 0, len(regular))
	for _, group := range regular {
		stageOptions = append(stageOptions, group.Stage)
	}

	return PostmortemsIndexView{Rows: rows, Filters: filters, Stages: stageOptions}
}

// postmortemsIndexAPIHandler serves the Astro postmortems-index page's
// data: the same buildPostmortemsIndex computation the old
// html/template handler used, as JSON.
func postmortemsIndexAPIHandler(cfg webcore.Config) http.HandlerFunc {
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

		view := buildPostmortemsIndex(cat, pf, r.URL.Query())
		webcore.WriteJSON(w, view)
	}
}

// postmortemAPIHandler serves the Astro postmortem-detail page's data:
// the same slug->path confinement and buildPostmortemPage computation
// the old html/template handler used, as JSON.
func postmortemAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pf, err := webcore.LoadProgress(cfg)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		path, ok := postmortemSlugs(pf)[r.PathValue("slug")]
		if !ok {
			webcore.WriteJSONError(w, http.StatusNotFound, "postmortem not found")
			return
		}

		cat, err := webcore.LoadCatalog(cfg)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		view, err := buildPostmortemPage(path, cat)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if idx, ok := attemptIndexForPostmortemPath(pf, path); ok {
			view.TicketID = webcore.TicketID(idx)
		}
		webcore.WriteJSON(w, view)
	}
}

// RegisterRoutes mounts the postmortem JSON API routes on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/postmortems", postmortemsIndexAPIHandler(cfg))
	mux.HandleFunc("GET /api/postmortems/{slug}", postmortemAPIHandler(cfg))
}
