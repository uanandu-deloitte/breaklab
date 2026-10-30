package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"breaklab/internal/docker"
	"breaklab/internal/runner"
)

func newRunner() *runner.Runner {
	r := runner.New(docker.New())
	if info, err := os.Stdin.Stat(); err == nil {
		r.Interactive = info.Mode()&os.ModeCharDevice != 0
	}
	return r
}

// parseFlags parses a command's flags; it returns done=true when the command
// should exit with code immediately (help requested or bad flags).
func parseFlags(fs *flag.FlagSet, args []string) (done bool, code int) {
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return true, 0
		}
		return true, 2
	}
	return false, 0
}

// report prints err and maps it to an exit code; Docker problems exit 2.
func report(err error) int {
	if errors.Is(err, docker.ErrUnavailable) {
		fmt.Fprintln(os.Stderr, "breaklab: Docker is not available. Start Docker and try again:")
		fmt.Fprintln(os.Stderr, "  macOS/Windows: open Docker Desktop")
		fmt.Fprintln(os.Stderr, "  Linux:         sudo systemctl start docker")
		fmt.Fprintf(os.Stderr, "(%v)\n", err)
		return 2
	}
	fmt.Fprintf(os.Stderr, "breaklab: %v\n", err)
	return 1
}

// commandContext is cancelled on Ctrl-C (SIGINT) or SIGTERM — the latter is
// what systemd sends on `systemctl stop`, so `breaklab serve` shuts down
// gracefully when run as a service too.
func commandContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func runStart(args []string) int {
	fs := flag.NewFlagSet("start", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: breaklab start [--yes] [--seed N] <lab>")
		fmt.Fprintln(os.Stderr, "\nBuilds a fresh container for the lab and breaks it. Prints the ticket.")
		fs.PrintDefaults()
	}
	yes := fs.Bool("yes", false, "stop a different active lab without asking")
	var seed int64
	fs.Int64Var(&seed, "seed", 0, "seed the variant selection for a reproducible pick")
	if done, code := parseFlags(fs, args); done {
		return code
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return 2
	}

	r := newRunner()
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			r.Seed = &seed
		}
	})

	ctx, cancel := commandContext()
	defer cancel()
	if err := r.Start(ctx, fs.Arg(0), *yes); err != nil {
		return report(err)
	}
	return 0
}

func runCheck(args []string) int {
	return runNoArgs("check", "Runs the lab's check and says whether the symptom is gone.", args, func(ctx context.Context, r *runner.Runner) int {
		res, err := r.Check(ctx)
		if err != nil {
			return report(err)
		}
		fmt.Println(res.Message)
		if !res.Passed {
			return 1
		}
		fmt.Printf("Variant was: %s\n", res.Variant)
		return 0
	})
}

func runHint(args []string) int {
	return runNoArgs("hint", "Reveals the next hint for the active lab.", args, func(ctx context.Context, r *runner.Runner) int {
		hint, err := r.Hint(ctx)
		if err != nil {
			return report(err)
		}
		fmt.Println(hint)
		return 0
	})
}

func runReset(args []string) int {
	return runNoArgs("reset", "Recreates the active lab's container with the same breakage.", args, func(ctx context.Context, r *runner.Runner) int {
		if err := r.Reset(ctx); err != nil {
			return report(err)
		}
		fmt.Println("Lab reset. Same problem, fresh container.")
		return 0
	})
}

func runStop(args []string) int {
	return runNoArgs("stop", "Removes the active lab's container.", args, func(ctx context.Context, r *runner.Runner) int {
		if err := r.Stop(ctx); err != nil {
			return report(err)
		}
		return 0
	})
}

func runNoArgs(name, desc string, args []string, fn func(context.Context, *runner.Runner) int) int {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: breaklab %s\n\n%s\n", name, desc)
	}
	fs.SetOutput(io.Discard)
	if done, code := parseFlags(fs, args); done {
		return code
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return 2
	}
	ctx, cancel := commandContext()
	defer cancel()
	return fn(ctx, newRunner())
}
