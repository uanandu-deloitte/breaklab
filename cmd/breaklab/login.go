package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// accountLinkPath is where `breaklab login` stores the hosted server
// URL and API token — CLI-local, gitignored state, same tier as
// .breaklab/session.json, never inside progress/progress.json.
var accountLinkPath = filepath.Join(".breaklab", "account.json")

// accountLink is the parsed shape of .breaklab/account.json.
type accountLink struct {
	Server string `json:"server"`
	Token  string `json:"token"`
}

// loadAccountLink reads path. A missing file is not an error — it
// returns (nil, nil), meaning "not linked to any account yet".
func loadAccountLink(path string) (*accountLink, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var link accountLink
	if err := json.Unmarshal(data, &link); err != nil {
		return nil, fmt.Errorf("%s is corrupt (%v)", path, err)
	}
	return &link, nil
}

// saveAccountLink writes link to path, mode 0600 since it carries a
// bearer token. A single os.WriteFile, not the temp-file-plus-rename
// dance progress.Save uses: unlike progress.json, this file is never
// written concurrently by more than one process.
func saveAccountLink(path string, link *accountLink) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(link, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// whoami asks server which account token belongs to, validating it in
// the process — a non-200 response means the token is invalid, revoked,
// or the server is unreachable.
func whoami(server, token string) (username string, err error) {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(server, "/")+"/sync/whoami", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("contacting %s: %w", server, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed: %s returned %s (check the token and try again)", server, resp.Status)
	}

	var payload struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("parsing response from %s: %w", server, err)
	}
	return payload.Username, nil
}

func runLogin(args []string) int {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: breaklab login --server <url> --token <token>")
		fmt.Fprintln(os.Stderr, "\nLinks this install to a hosted account for `breaklab sync`.")
		fmt.Fprintln(os.Stderr, "Generate a token on that account's settings page first.")
		fs.PrintDefaults()
	}
	server := fs.String("server", "", "hosted server URL (required)")
	token := fs.String("token", "", "API token from the account settings page (required)")
	if done, code := parseFlags(fs, args); done {
		return code
	}
	if *server == "" || *token == "" {
		fs.Usage()
		return 2
	}

	username, err := whoami(*server, *token)
	if err != nil {
		return report(err)
	}

	if err := saveAccountLink(accountLinkPath, &accountLink{Server: *server, Token: *token}); err != nil {
		return report(err)
	}
	fmt.Printf("Linked to account %s on %s\n", username, *server)
	return 0
}
