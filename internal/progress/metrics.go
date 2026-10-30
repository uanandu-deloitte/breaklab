package progress

import (
	"fmt"
	"sort"
	"time"

	"breaklab/internal/catalog"
)

// solveTimes returns the resolution time (SolvedAt - StartedAt) of every
// solved or complete attempt that recorded a SolvedAt.
func solveTimes(attempts []Attempt) []time.Duration {
	var out []time.Duration
	for _, a := range attempts {
		if (a.Status != StatusSolved && a.Status != StatusComplete) || a.SolvedAt == nil {
			continue
		}
		out = append(out, a.SolvedAt.Sub(a.StartedAt))
	}
	return out
}

func medianDuration(ds []time.Duration) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	sorted := append([]time.Duration(nil), ds...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}

// AvgResolutionTime is the mean resolution time across every
// solved/complete attempt, for the dashboard's stat tile. Zero when
// there's no qualifying attempt yet.
func AvgResolutionTime(attempts []Attempt) time.Duration {
	times := solveTimes(attempts)
	if len(times) == 0 {
		return 0
	}
	var total time.Duration
	for _, t := range times {
		total += t
	}
	return total / time.Duration(len(times))
}

// FormatDuration renders d as a short "2h45m"/"18m45s"/"12s" label for
// display. Zero renders as "".
func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return ""
	}
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	switch {
	case h > 0:
		return fmt.Sprintf("%dh%dm", h, m)
	case m > 0:
		return fmt.Sprintf("%dm%ds", m, s)
	default:
		return fmt.Sprintf("%ds", s)
	}
}

// PostmortemCounts splits attempts into drafts (solved, write-up still
// pending) and published (complete, with a validated postmortem on
// record) for the dashboard's postmortem-records stat tile.
func PostmortemCounts(attempts []Attempt) (drafts, published int) {
	for _, a := range attempts {
		switch a.Status {
		case StatusComplete:
			published++
		case StatusSolved:
			drafts++
		}
	}
	return drafts, published
}

// DifficultyBucket is one difficulty level's share of the catalog: how
// many labs exist at that difficulty, and how many of those are solved.
type DifficultyBucket struct {
	Difficulty int
	Count      int
	Solved     int
}

// TrainingVelocity is the challenge browser's aggregate progress
// summary, computed entirely from the catalog and the user's own
// recorded attempts — no new persisted data.
type TrainingVelocity struct {
	PercentSolved   float64
	SolvedCount     int
	TotalCount      int
	DifficultyCurve []DifficultyBucket
	MedianSolveTime time.Duration

	// MedianSolveTimeDisplay is MedianSolveTime pre-formatted for
	// display (its native Stringer output, e.g. "18m45.000000001s", is
	// not the short form the rest of the UI uses).
	MedianSolveTimeDisplay string
}

// BuildTrainingVelocity computes TrainingVelocity for cat/pf.
func BuildTrainingVelocity(cat *catalog.Catalog, pf *ProgressFile) TrainingVelocity {
	latest := LatestCompleteAttempts(pf)

	buckets := make(map[int]*DifficultyBucket)
	total := 0
	for _, group := range cat.RegularStages() {
		for _, lab := range group.Labs {
			total++
			b, ok := buckets[lab.Difficulty]
			if !ok {
				b = &DifficultyBucket{Difficulty: lab.Difficulty}
				buckets[lab.Difficulty] = b
			}
			b.Count++
			if _, solved := latest[lab.ID]; solved {
				b.Solved++
			}
		}
	}

	difficulties := make([]int, 0, len(buckets))
	for d := range buckets {
		difficulties = append(difficulties, d)
	}
	sort.Ints(difficulties)
	curve := make([]DifficultyBucket, 0, len(difficulties))
	for _, d := range difficulties {
		curve = append(curve, *buckets[d])
	}

	var percent float64
	if total > 0 {
		percent = 100 * float64(len(latest)) / float64(total)
	}

	median := medianDuration(solveTimes(pf.Attempts))

	return TrainingVelocity{
		PercentSolved:          percent,
		SolvedCount:            len(latest),
		TotalCount:             total,
		DifficultyCurve:        curve,
		MedianSolveTime:        median,
		MedianSolveTimeDisplay: FormatDuration(median),
	}
}
