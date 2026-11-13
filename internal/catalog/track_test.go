package catalog

import "testing"

func stageGroup(id, track string) StageGroup {
	return StageGroup{Stage: Stage{ID: id, Track: track}}
}

func TestGroupByTrackAllSoloStagesStayUngrouped(t *testing.T) {
	groups := []StageGroup{stageGroup("linux", ""), stageGroup("git", ""), stageGroup("docker", "")}

	entries := GroupByTrack(groups)
	if len(entries) != 3 {
		t.Fatalf("len(entries) = %d, want 3 (one per solo stage)", len(entries))
	}
	for i, id := range []string{"linux", "git", "docker"} {
		e := entries[i]
		if e.Track != "" || e.Stage == nil || e.Stage.Stage.ID != id {
			t.Errorf("entries[%d] = %+v, want a solo entry for %q", i, e, id)
		}
	}
}

func TestGroupByTrackCollapsesAdjacentRunIntoOneEntry(t *testing.T) {
	groups := []StageGroup{
		stageGroup("k8s-fundamentals", ""),
		stageGroup("k8s-developer", "k8s-specialization"),
		stageGroup("k8s-admin", "k8s-specialization"),
		stageGroup("k8s-security", "k8s-specialization"),
		stageGroup("monitoring", ""),
	}

	entries := GroupByTrack(groups)
	if len(entries) != 3 {
		t.Fatalf("len(entries) = %d, want 3 (solo, one grouped run, solo)", len(entries))
	}
	if entries[0].Track != "" || entries[0].Stage.Stage.ID != "k8s-fundamentals" {
		t.Errorf("entries[0] = %+v, want the solo k8s-fundamentals stage", entries[0])
	}
	grouped := entries[1]
	if grouped.Track != "k8s-specialization" || len(grouped.Stages) != 3 {
		t.Fatalf("entries[1] = %+v, want a 3-member k8s-specialization group", grouped)
	}
	for i, id := range []string{"k8s-developer", "k8s-admin", "k8s-security"} {
		if grouped.Stages[i].Stage.ID != id {
			t.Errorf("grouped.Stages[%d].ID = %q, want %q (roadmap.yaml order preserved)", i, grouped.Stages[i].Stage.ID, id)
		}
	}
	if entries[2].Track != "" || entries[2].Stage.Stage.ID != "monitoring" {
		t.Errorf("entries[2] = %+v, want the solo monitoring stage", entries[2])
	}
}

// Two non-adjacent runs of the same Track value must stay as two
// separate groups, matching track.go's doc comment: "Member stages
// must be adjacent in roadmap.yaml; GroupByTrack groups them by
// scanning for runs of equal Track, not by matching the value
// globally."
func TestGroupByTrackNonAdjacentSameTrackStaysTwoGroups(t *testing.T) {
	groups := []StageGroup{
		stageGroup("a1", "x"),
		stageGroup("a2", "x"),
		stageGroup("between", ""),
		stageGroup("b1", "x"),
	}

	entries := GroupByTrack(groups)
	if len(entries) != 3 {
		t.Fatalf("len(entries) = %d, want 3 (group, solo, group)", len(entries))
	}
	if len(entries[0].Stages) != 2 {
		t.Errorf("entries[0] = %+v, want the first 2-member run", entries[0])
	}
	if entries[1].Track != "" || entries[1].Stage.Stage.ID != "between" {
		t.Errorf("entries[1] = %+v, want the solo stage splitting the two runs", entries[1])
	}
	if len(entries[2].Stages) != 1 || entries[2].Stages[0].Stage.ID != "b1" {
		t.Errorf("entries[2] = %+v, want a separate 1-member run for b1", entries[2])
	}
}

func TestGroupByTrackEmptyInputReturnsEmpty(t *testing.T) {
	entries := GroupByTrack(nil)
	if len(entries) != 0 {
		t.Errorf("entries = %+v, want empty for no stages", entries)
	}
}

func TestGroupByTrackWholeRoadmapIsOneTrack(t *testing.T) {
	groups := []StageGroup{stageGroup("a", "x"), stageGroup("b", "x")}

	entries := GroupByTrack(groups)
	if len(entries) != 1 || len(entries[0].Stages) != 2 {
		t.Fatalf("entries = %+v, want a single 2-member group", entries)
	}
}
