package progress

import (
	"reflect"
	"testing"

	"breaklab/internal/catalog"
)

func TestStageCompletion(t *testing.T) {
	cat := &catalog.Catalog{
		Stages: []catalog.StageGroup{
			{
				Stage: catalog.Stage{ID: "linux", Title: "Linux fundamentals"},
				Labs: []catalog.Lab{
					{ID: "web-down"},
					{ID: "disk-full"},
				},
			},
			{
				Stage: catalog.Stage{ID: "git", Title: "Git"},
				Labs: []catalog.Lab{
					{ID: "git-lost-work"},
				},
			},
		},
	}

	pf := &ProgressFile{Attempts: []Attempt{
		{Lab: "web-down", Status: StatusSolved},
		{Lab: "disk-full", Status: StatusActive},
		{Lab: "git-lost-work", Status: StatusAbandoned},
	}}

	got := StageCompletion(cat, pf)
	want := []StageStatus{
		{ID: "linux", Title: "Linux fundamentals", Completed: 1, Total: 2},
		{ID: "git", Title: "Git", Completed: 0, Total: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestStageCompletionCountsCompleteStatusToo(t *testing.T) {
	cat := &catalog.Catalog{Stages: []catalog.StageGroup{
		{Stage: catalog.Stage{ID: "linux"}, Labs: []catalog.Lab{{ID: "web-down"}}},
	}}
	pf := &ProgressFile{Attempts: []Attempt{{Lab: "web-down", Status: StatusComplete}}}

	got := StageCompletion(cat, pf)
	if len(got) != 1 || got[0].Completed != 1 {
		t.Errorf("got %+v", got)
	}
}
