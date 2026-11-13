package main

import (
	"fmt"
	"os"
	"path/filepath"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

func runList(args []string) int {
	cat, err := catalog.LoadCatalog("labs", "roadmap.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "breaklab: %v\n", err)
		return 1
	}

	pf, err := progress.Load(filepath.Join("progress", "progress.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "breaklab: %v\n", err)
		return 1
	}
	progress.ApplyStatus(cat, pf)

	for _, warning := range cat.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}

	for _, entry := range catalog.GroupByTrack(cat.RegularStages()) {
		if entry.Stage != nil {
			printStageLabs(*entry.Stage, "")
			continue
		}
		fmt.Printf("%s — pick one track:\n", entry.Track)
		for _, group := range entry.Stages {
			printStageLabs(group, "  ")
		}
	}

	if capstone := progress.Capstone(cat, pf); capstone != nil {
		fmt.Println()
		if capstone.Unlocked {
			fmt.Printf("%s (unlocked)\n", capstone.Title)
			for _, lab := range cat.CapstoneStage().Labs {
				fmt.Printf("  %-20s %-40s difficulty=%d  %s\n", lab.ID, lab.Title, lab.Difficulty, lab.Status)
			}
		} else {
			fmt.Printf("%s — locked until every stage above is complete\n", capstone.Title)
		}
	}

	return 0
}

// printStageLabs prints one stage's title and labs, each line prefixed
// by indent (used to nest a track group's member stages under its
// shared "pick one track" header). Mirrors runList's old inline loop,
// skipping a stage with no labs rather than printing an empty heading.
func printStageLabs(group catalog.StageGroup, indent string) {
	if len(group.Labs) == 0 {
		return
	}
	fmt.Printf("%s%s\n", indent, group.Stage.Title)
	for _, lab := range group.Labs {
		fmt.Printf("%s  %-20s %-40s difficulty=%d  %s\n", indent, lab.ID, lab.Title, lab.Difficulty, lab.Status)
	}
}
