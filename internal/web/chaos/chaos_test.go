package chaos

import (
	"testing"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

func TestBuildChaosGroupsByTagAcrossStages(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{
			{ID: "zombie-pileup", Title: "Zombie pileup", Difficulty: 2, Tags: []string{"linux", "chaos-engineering", "fault-injection"}},
		}},
		{Stage: catalog.Stage{ID: "k8s-developer"}, Labs: []catalog.Lab{
			{ID: "oom-killed", Title: "OOM killed", Difficulty: 1, Tags: []string{"kubernetes", "chaos-engineering", "fault-injection"}},
			{ID: "selector-mismatch", Title: "Selector mismatch", Difficulty: 2, Tags: []string{"kubernetes", "chaos-engineering", "blast-radius"}},
			{ID: "untagged-lab", Title: "Untagged", Difficulty: 1, Tags: []string{"kubernetes"}},
		}},
	}}

	view := buildChaos(cat, &progress.ProgressFile{})

	var faultInjection, blastRadius ChaosConceptView
	for _, c := range view.Concepts {
		switch c.ID {
		case "fault-injection":
			faultInjection = c
		case "blast-radius":
			blastRadius = c
		}
	}

	if len(faultInjection.Labs) != 2 || faultInjection.Labs[0].ID != "oom-killed" || faultInjection.Labs[1].ID != "zombie-pileup" {
		t.Fatalf("fault-injection labs = %+v, want [oom-killed, zombie-pileup] ordered by difficulty", faultInjection.Labs)
	}
	if len(blastRadius.Labs) != 1 || blastRadius.Labs[0].ID != "selector-mismatch" {
		t.Fatalf("blast-radius labs = %+v, want [selector-mismatch]", blastRadius.Labs)
	}
}
