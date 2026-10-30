package main

import (
	"flag"
	"fmt"
	"os"
)

type command struct {
	summary string
	run     func(args []string) int
}

var commands = map[string]command{
	"list": {
		summary: "List labs grouped by roadmap stage",
		run:     runList,
	},
	"start": {
		summary: "Start a lab: build a container and break it",
		run:     runStart,
	},
	"check": {
		summary: "Check whether the active lab is fixed",
		run:     runCheck,
	},
	"hint": {
		summary: "Reveal the next hint for the active lab",
		run:     runHint,
	},
	"reset": {
		summary: "Recreate the active lab with the same breakage",
		run:     runReset,
	},
	"stop": {
		summary: "Remove the active lab's container",
		run:     runStop,
	},
	"status": {
		summary: "Show per-stage progress, the active lab, and your streak",
		run:     runStatus,
	},
	"postmortem": {
		summary: "Write up or validate a postmortem for a solved lab",
		run:     runPostmortem,
	},
	"publish": {
		summary: "Rebuild site/ from progress, labs, and published postmortems",
		run:     runPublish,
	},
	"serve": {
		summary: "Serve the local dashboard, roadmap, and challenge browser",
		run:     runServe,
	},
	"login": {
		summary: "Link this install to a hosted account for `sync`",
		run:     runLogin,
	},
	"sync": {
		summary: "Push/pull progress with the linked hosted account",
		run:     runSync,
	},
}

var commandOrder = []string{"list", "start", "check", "hint", "reset", "stop", "status", "postmortem", "publish", "serve", "login", "sync"}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flag.Usage = printUsage
	if len(args) == 0 {
		printUsage()
		return 2
	}

	if args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return 0
	}

	cmd, ok := commands[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "breaklab: unknown command %q\n\n", args[0])
		printUsage()
		return 2
	}

	return cmd.run(args[1:])
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: breaklab <command> [arguments]")
	fmt.Fprintln(os.Stderr, "\nCommands:")
	for _, name := range commandOrder {
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", name, commands[name].summary)
	}
}
