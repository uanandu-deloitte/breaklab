package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"breaklab/internal/progress"
)

// progressPath is where `breaklab sync` reads and writes local
// progress — the same path `breaklab serve`'s Config.ProgressPath
// points at.
var syncProgressPath = filepath.Join("progress", "progress.json")

// syncAttemptsPayload mirrors internal/web's wire shape for
// /sync/attempts. Duplicated rather than imported: this is the CLI's
// HTTP client, that package's HTTP server, and sharing one tiny JSON
// struct isn't worth coupling the CLI to the web package.
type syncAttemptsPayload struct {
	Attempts []progress.Attempt `json:"attempts"`
}

// isSyncable reports whether status is one `breaklab sync` ever moves.
// An "active" attempt only makes sense on the machine actually running
// the lab, so it never crosses the wire (design.md "Sync only ever
// moves terminal attempts").
func isSyncable(status string) bool {
	switch status {
	case progress.StatusSolved, progress.StatusComplete, progress.StatusAbandoned, progress.StatusSkipped:
		return true
	default:
		return false
	}
}

func postAttempts(link *accountLink, attempts []progress.Attempt) error {
	body, err := json.Marshal(syncAttemptsPayload{Attempts: attempts})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(link.Server, "/")+"/sync/attempts", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+link.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("pushing to %s: %w", link.Server, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("pushing to %s: server returned %s", link.Server, resp.Status)
	}
	return nil
}

func getAttempts(link *accountLink) ([]progress.Attempt, error) {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(link.Server, "/")+"/sync/attempts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+link.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pulling from %s: %w", link.Server, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pulling from %s: server returned %s", link.Server, resp.Status)
	}

	var payload syncAttemptsPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("parsing response from %s: %w", link.Server, err)
	}
	return payload.Attempts, nil
}

// syncOnce pushes every syncable local attempt not yet known to the
// hosted account, then pulls every hosted attempt not yet present
// locally, merging by Attempt.ID. It never touches Docker or any lab
// execution state — only progressPath and the HTTP calls above.
func syncOnce(progressPath string, link *accountLink) (pushed, pulled int, err error) {
	if err := progress.BackfillAttemptIDs(progressPath); err != nil {
		return 0, 0, err
	}
	pf, err := progress.Load(progressPath)
	if err != nil {
		return 0, 0, err
	}

	var toPush []progress.Attempt
	for _, a := range pf.Attempts {
		if isSyncable(a.Status) {
			toPush = append(toPush, a)
		}
	}
	if len(toPush) > 0 {
		if err := postAttempts(link, toPush); err != nil {
			return 0, 0, err
		}
	}

	hosted, err := getAttempts(link)
	if err != nil {
		return 0, 0, err
	}

	known := make(map[string]bool, len(pf.Attempts))
	for _, a := range pf.Attempts {
		if a.ID != "" {
			known[a.ID] = true
		}
	}
	for _, a := range hosted {
		if a.ID == "" || known[a.ID] {
			continue
		}
		pf.Attempts = append(pf.Attempts, a)
		pulled++
	}
	if pulled > 0 {
		if err := progress.Save(progressPath, pf); err != nil {
			return 0, 0, err
		}
	}
	return len(toPush), pulled, nil
}

func runSync(args []string) int {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: breaklab sync")
		fmt.Fprintln(os.Stderr, "\nPushes solved/complete/abandoned/skipped progress to the linked hosted")
		fmt.Fprintln(os.Stderr, "account and pulls down progress synced from other devices.")
		fmt.Fprintln(os.Stderr, "Run `breaklab login --server <url> --token <token>` first.")
	}
	if done, code := parseFlags(fs, args); done {
		return code
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return 2
	}

	link, err := loadAccountLink(accountLinkPath)
	if err != nil {
		return report(err)
	}
	if link == nil {
		fmt.Fprintln(os.Stderr, "breaklab: not logged in. Run `breaklab login --server <url> --token <token>` first.")
		return 1
	}

	pushed, pulled, err := syncOnce(syncProgressPath, link)
	if err != nil {
		return report(err)
	}
	fmt.Printf("Synced: pushed %d, pulled %d\n", pushed, pulled)
	return 0
}
