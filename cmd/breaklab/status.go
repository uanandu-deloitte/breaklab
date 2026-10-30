package main

import (
	"context"
	"fmt"
	"strings"

	"breaklab/internal/runner"
)

func runStatus(args []string) int {
	return runNoArgs("status", "Shows per-stage progress, the active lab, and your streak.", args, func(ctx context.Context, r *runner.Runner) int {
		view, err := r.Status(ctx)
		if err != nil {
			return report(err)
		}
		printStatus(view)
		return 0
	})
}

func printStatus(view runner.StatusView) {
	for _, s := range view.Stages {
		fmt.Printf("%-24s %d/%d labs complete\n", s.Title, s.Completed, s.Total)
	}

	if view.Capstone != nil {
		if view.Capstone.Unlocked {
			fmt.Printf("%-24s %d/%d labs complete\n", view.Capstone.Title, view.Capstone.Completed, view.Capstone.Total)
		} else {
			fmt.Printf("%-24s locked (finish every stage above first)\n", view.Capstone.Title)
		}
	}

	fmt.Println()
	if view.Active != nil {
		fmt.Printf("Active: %s (elapsed %dm, %d hint(s) used)\n", view.Active.Lab, view.Active.ElapsedMinutes, view.Active.HintsUsed)
	} else {
		fmt.Println("No lab active.")
	}
	fmt.Printf("Streak: %d day(s)\n", view.Streak)

	if len(view.ReviewDue) > 0 {
		fmt.Printf("Review due: %s\n", strings.Join(view.ReviewDue, ", "))
	}

	if len(view.Patterns) > 0 {
		fmt.Println("\nPatterns:")
		for _, p := range view.Patterns {
			fmt.Printf("  %-16s avg %.1f hints over %d labs\n", p.Tag, p.AverageHintsUsed, p.Count)
		}
	}
}
