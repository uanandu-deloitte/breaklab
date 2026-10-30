package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"breaklab/internal/site"
)

func runPublish(args []string) int {
	fs := flag.NewFlagSet("publish", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: breaklab publish")
		fmt.Fprintln(os.Stderr, "\nRebuilds site/ from progress.json, roadmap.yaml, lab metadata, and published postmortems.")
	}
	if done, code := parseFlags(fs, args); done {
		return code
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return 2
	}

	cfg := site.Config{
		LabsDir:        "labs",
		RoadmapPath:    "roadmap.yaml",
		ProgressPath:   filepath.Join("progress", "progress.json"),
		TemplatesDir:   filepath.Join("templates", "site"),
		PostmortemsDir: "postmortems",
		OutDir:         "site",
	}
	if err := site.Build(cfg); err != nil {
		return report(err)
	}
	fmt.Println("Published site/")
	return 0
}
