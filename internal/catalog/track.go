package catalog

// RoadmapEntry is one slot in a regular-stage sequence: either a solo
// stage, or a run of stages sharing a Track — alternatives a learner
// picks one of, not a required sequence of all of them.
type RoadmapEntry struct {
	// Track is non-empty only for a grouped entry; Stages then holds
	// every member stage, in roadmap.yaml order. A solo entry leaves
	// this empty and sets Stage instead.
	Track  string
	Stage  *StageGroup
	Stages []StageGroup
}

// GroupByTrack scans groups (already in roadmap.yaml order) for runs of
// adjacent stages sharing the same non-empty Track and collapses each
// run into one grouped RoadmapEntry. A stage with an empty Track, or
// whose Track differs from its predecessor's, stays its own solo entry.
func GroupByTrack(groups []StageGroup) []RoadmapEntry {
	entries := make([]RoadmapEntry, 0, len(groups))
	for i := 0; i < len(groups); i++ {
		track := groups[i].Stage.Track
		if track == "" {
			g := groups[i]
			entries = append(entries, RoadmapEntry{Stage: &g})
			continue
		}
		j := i + 1
		for j < len(groups) && groups[j].Stage.Track == track {
			j++
		}
		run := append([]StageGroup(nil), groups[i:j]...)
		entries = append(entries, RoadmapEntry{Track: track, Stages: run})
		i = j - 1
	}
	return entries
}
