package main

import (
	"flag"
	"fmt"
	"os"
)

func runPostmortem(args []string) int {
	fs := flag.NewFlagSet("postmortem", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: breaklab postmortem [--validate] [path]")
		fmt.Fprintln(os.Stderr, "\nCreates a postmortem for the most recently solved lab and opens it in $EDITOR.")
		fmt.Fprintln(os.Stderr, "With --validate, checks a postmortem's required sections instead of creating one.")
		fs.PrintDefaults()
	}
	validate := fs.Bool("validate", false, "validate a postmortem's required sections instead of creating one")
	if done, code := parseFlags(fs, args); done {
		return code
	}

	ctx, cancel := commandContext()
	defer cancel()
	r := newRunner()

	if !*validate {
		if fs.NArg() != 0 {
			fs.Usage()
			return 2
		}
		path, err := r.Postmortem(ctx)
		if err != nil {
			return report(err)
		}
		fmt.Printf("Created %s\n", path)
		return 0
	}

	if fs.NArg() > 1 {
		fs.Usage()
		return 2
	}
	var path string
	if fs.NArg() == 1 {
		path = fs.Arg(0)
	}

	res, err := r.ValidatePostmortem(ctx, path)
	if err != nil {
		return report(err)
	}
	if res.RootCauseWarning != "" {
		fmt.Println("Warning:", res.RootCauseWarning)
	}
	if !res.Valid {
		fmt.Println("Postmortem is incomplete. Needs more detail in:")
		for _, name := range res.Failing {
			fmt.Printf("  - %s\n", name)
		}
		return 1
	}
	fmt.Println("Postmortem looks complete.")
	return 0
}
