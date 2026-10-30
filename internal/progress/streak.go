package progress

import "time"

// Streak returns the current daily streak as of now: the number of
// consecutive local-calendar days, counting backward, with at least one
// solved/complete attempt (design.md Decision 7). If today has no
// solved attempt yet, that alone doesn't break the streak — the day
// isn't over — so counting starts from yesterday instead; from there it
// stops at the first day with no activity.
func Streak(pf *ProgressFile, now time.Time) int {
	activeDays := make(map[string]bool, len(pf.Attempts))
	for _, a := range pf.Attempts {
		if (a.Status == StatusSolved || a.Status == StatusComplete) && a.SolvedAt != nil {
			activeDays[dateKey(a.SolvedAt.Local())] = true
		}
	}

	day := now.Local()
	if !activeDays[dateKey(day)] {
		day = day.AddDate(0, 0, -1)
	}

	count := 0
	for activeDays[dateKey(day)] {
		count++
		day = day.AddDate(0, 0, -1)
	}
	return count
}

func dateKey(t time.Time) string {
	return t.Format("2006-01-02")
}
