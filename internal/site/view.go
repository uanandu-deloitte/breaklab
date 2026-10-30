// Package site builds the static, offline-capable site/ directory from
// committed data only: roadmap.yaml, lab metadata, progress/progress.json,
// and postmortems/*.md marked publish: true.
package site

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

// CompletedLabView is one completed lab's row under its stage.
type CompletedLabView struct {
	ID             string
	Title          string
	Date           string
	Difficulty     int
	MinutesToFix   int
	PostmortemHref string // empty unless its postmortem is published
}

// StageView is one roadmap stage's progress bar plus its completed labs.
type StageView struct {
	ID        string
	Title     string
	Completed int
	Total     int
	Percent   int
	Labs      []CompletedLabView
}

// Totals is the summary shown at the bottom of the site.
type Totals struct {
	LabsCompleted int
	VariantsSeen  int
	Streak        int
}

// IndexView is everything site/index.html renders.
type IndexView struct {
	Stages []StageView
	Totals Totals
}

// buildStageView builds one stage's progress bar and completed-lab rows.
func buildStageView(group catalog.StageGroup, latest map[string]progress.Attempt, published map[string]bool) StageView {
	sv := StageView{ID: group.Stage.ID, Title: group.Stage.Title, Total: len(group.Labs)}
	for _, lab := range group.Labs {
		attempt, ok := latest[lab.ID]
		if !ok {
			continue
		}
		sv.Completed++

		row := CompletedLabView{
			ID:           lab.ID,
			Title:        lab.Title,
			Date:         attempt.SolvedAt.Local().Format("2006-01-02"),
			Difficulty:   lab.Difficulty,
			MinutesToFix: int(attempt.SolvedAt.Sub(attempt.StartedAt).Minutes()),
		}
		if published[attempt.Postmortem] {
			row.PostmortemHref = "postmortems/" + slug(attempt.Postmortem) + ".html"
		}
		sv.Labs = append(sv.Labs, row)
	}
	if sv.Total > 0 {
		sv.Percent = sv.Completed * 100 / sv.Total
	}
	return sv
}

// BuildIndex assembles IndexView from the catalog and progress.json. A
// lab counts as completed only when it has an attempt with status
// "complete" (design.md Decision 2); its date/difficulty/minutes-to-fix
// always show, but PostmortemHref is set only when published[path] is
// true (design.md Decision 1). VariantsSeen counts distinct lab+variant
// pairs across every attempt, any status (design.md Decision 3).
func BuildIndex(cat *catalog.Catalog, pf *progress.ProgressFile, published map[string]bool) IndexView {
	latest := progress.LatestCompleteAttempts(pf)
	regular := cat.RegularStages()
	stages := make([]StageView, 0, len(regular)+1)

	for _, group := range regular {
		stages = append(stages, buildStageView(group, latest, published))
	}

	// The capstone is a surprise, earned finale: it stays off the public
	// site entirely until there's something to show for it, rather than
	// appearing as a visibly locked stage before it's been completed.
	if capstone := cat.CapstoneStage(); capstone != nil {
		sv := buildStageView(*capstone, latest, published)
		if sv.Completed > 0 {
			stages = append(stages, sv)
		}
	}

	seenVariants := make(map[string]bool)
	for _, a := range pf.Attempts {
		seenVariants[a.Lab+"\x00"+a.Variant] = true
	}

	return IndexView{
		Stages: stages,
		Totals: Totals{
			LabsCompleted: len(latest),
			VariantsSeen:  len(seenVariants),
			Streak:        progress.Streak(pf, time.Now()),
		},
	}
}

// slug derives a postmortem's page basename (no extension) from its
// file path.
func slug(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// PublishedPostmortemRow is one entry in the flat public postmortems
// index — every published write-up, independent of which stage it's
// nested under on the main progress index.
type PublishedPostmortemRow struct {
	Title        string
	StageTitle   string
	Date         string
	Difficulty   int
	MinutesToFix int
	Href         string // "<slug>.html", relative to postmortems/
}

// PostmortemsIndexView is everything site/postmortems/index.html renders.
type PostmortemsIndexView struct {
	Postmortems []PublishedPostmortemRow
}

// BuildPostmortemsIndex lists every published postmortem across every
// stage (regular and capstone), most recently solved first, from the
// same inputs BuildIndex already takes — no new data source. A lab
// with no complete attempt, or whose postmortem isn't in published, is
// skipped: the "nothing leaks" rule this function must uphold is the
// same one governing every other page in site/.
func BuildPostmortemsIndex(cat *catalog.Catalog, pf *progress.ProgressFile, published map[string]bool) PostmortemsIndexView {
	latest := progress.LatestCompleteAttempts(pf)

	groups := cat.RegularStages()
	if capstone := cat.CapstoneStage(); capstone != nil {
		groups = append(append([]catalog.StageGroup{}, groups...), *capstone)
	}

	var rows []PublishedPostmortemRow
	for _, group := range groups {
		for _, lab := range group.Labs {
			attempt, ok := latest[lab.ID]
			if !ok || !published[attempt.Postmortem] {
				continue
			}
			rows = append(rows, PublishedPostmortemRow{
				Title:        lab.Title,
				StageTitle:   group.Stage.Title,
				Date:         attempt.SolvedAt.Local().Format("2006-01-02"),
				Difficulty:   lab.Difficulty,
				MinutesToFix: int(attempt.SolvedAt.Sub(attempt.StartedAt).Minutes()),
				Href:         slug(attempt.Postmortem) + ".html",
			})
		}
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Date > rows[j].Date })

	return PostmortemsIndexView{Postmortems: rows}
}
