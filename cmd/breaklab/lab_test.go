package main

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"breaklab/internal/docker"
)

func TestReportMapsDockerUnavailableToExit2(t *testing.T) {
	if code := report(fmt.Errorf("wrapped: %w", docker.ErrUnavailable)); code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if code := report(fmt.Errorf("something else")); code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}

func TestLabCommandsHelp(t *testing.T) {
	for _, name := range []string{"start", "check", "hint", "reset", "stop", "status", "postmortem", "publish", "serve"} {
		if code := run([]string{name, "--help"}); code != 0 {
			t.Errorf("%s --help exit = %d, want 0", name, code)
		}
	}
	if code := run([]string{"start"}); code != 2 {
		t.Errorf("start without a lab exit = %d, want 2", code)
	}
}

func TestHintRegisteredAsCommand(t *testing.T) {
	if _, ok := commands["hint"]; !ok {
		t.Error(`commands["hint"] not registered`)
	}
	if !slices.Contains(commandOrder, "hint") {
		t.Errorf("commandOrder = %v, missing %q", commandOrder, "hint")
	}
}

func TestStatusRegisteredAsCommand(t *testing.T) {
	if _, ok := commands["status"]; !ok {
		t.Error(`commands["status"] not registered`)
	}
	if !slices.Contains(commandOrder, "status") {
		t.Errorf("commandOrder = %v, missing %q", commandOrder, "status")
	}
}

func TestPostmortemRegisteredAsCommand(t *testing.T) {
	if _, ok := commands["postmortem"]; !ok {
		t.Error(`commands["postmortem"] not registered`)
	}
	if !slices.Contains(commandOrder, "postmortem") {
		t.Errorf("commandOrder = %v, missing %q", commandOrder, "postmortem")
	}
}

func TestPublishRegisteredAsCommand(t *testing.T) {
	if _, ok := commands["publish"]; !ok {
		t.Error(`commands["publish"] not registered`)
	}
	if !slices.Contains(commandOrder, "publish") {
		t.Errorf("commandOrder = %v, missing %q", commandOrder, "publish")
	}
}

func TestServeRegisteredAsCommand(t *testing.T) {
	if _, ok := commands["serve"]; !ok {
		t.Error(`commands["serve"] not registered`)
	}
	if !slices.Contains(commandOrder, "serve") {
		t.Errorf("commandOrder = %v, missing %q", commandOrder, "serve")
	}
}

func TestServeDefaultsToLoopback(t *testing.T) {
	fs, host, port, db, frontendDir := serveFlags()
	if err := fs.Parse(nil); err != nil {
		t.Fatal(err)
	}
	if *host != "127.0.0.1" {
		t.Errorf("default --host = %q, want 127.0.0.1", *host)
	}
	if *port != 8787 {
		t.Errorf("default --port = %d, want 8787", *port)
	}
	if *db != "" {
		t.Errorf("default --db = %q, want empty (accounts are opt-in)", *db)
	}
	if want := filepath.Join("frontend", "dist"); *frontendDir != want {
		t.Errorf("default --frontend = %q, want %q (migrated pages have no template fallback)", *frontendDir, want)
	}

	fs, host, _, _, _ = serveFlags()
	if err := fs.Parse([]string{"--host", "0.0.0.0"}); err != nil {
		t.Fatal(err)
	}
	if *host != "0.0.0.0" {
		t.Errorf("explicit --host = %q, want 0.0.0.0", *host)
	}
}
