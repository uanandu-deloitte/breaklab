package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"breaklab/internal/progress"
)

// withAccountLinkPath points the package-level accountLinkPath at path
// for the duration of the test, restoring it afterward — runLogin/runSync
// read that var rather than taking a path parameter, so tests isolate it
// the same way they isolate any other shared package state.
func withAccountLinkPath(t *testing.T, path string) {
	t.Helper()
	old := accountLinkPath
	accountLinkPath = path
	t.Cleanup(func() { accountLinkPath = old })
}

// withSyncProgressPath isolates runSync's progress.json path the same
// way withAccountLinkPath isolates its account-link path.
func withSyncProgressPath(t *testing.T, path string) {
	t.Helper()
	old := syncProgressPath
	syncProgressPath = path
	t.Cleanup(func() { syncProgressPath = old })
}

func validAuth(r *http.Request, token string) bool {
	return r.Header.Get("Authorization") == "Bearer "+token
}

// newFakeAccountServer serves the exact wire contract internal/web's
// /sync/* routes expose, without depending on that package — this is a
// test double for the CLI's HTTP client, not a copy of the real
// server. Its store is deduplicated by attempt ID on push, mirroring
// the real server's `INSERT OR IGNORE`, so round-trip tests see the
// same idempotent behavior the real one provides.
func newFakeAccountServer(t *testing.T, validToken string, seed []progress.Attempt) (*httptest.Server, *[]progress.Attempt) {
	t.Helper()
	store := append([]progress.Attempt(nil), seed...)
	known := make(map[string]bool, len(store))
	for _, a := range store {
		known[a.ID] = true
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /sync/whoami", func(w http.ResponseWriter, r *http.Request) {
		if !validAuth(r, validToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"username": "ada"})
	})
	mux.HandleFunc("GET /sync/attempts", func(w http.ResponseWriter, r *http.Request) {
		if !validAuth(r, validToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(syncAttemptsPayload{Attempts: store})
	})
	mux.HandleFunc("POST /sync/attempts", func(w http.ResponseWriter, r *http.Request) {
		if !validAuth(r, validToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var payload syncAttemptsPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, a := range payload.Attempts {
			if known[a.ID] {
				continue
			}
			known[a.ID] = true
			store = append(store, a)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, &store
}

func TestWhoamiValidToken(t *testing.T) {
	srv, _ := newFakeAccountServer(t, "good-token", nil)

	username, err := whoami(srv.URL, "good-token")
	if err != nil || username != "ada" {
		t.Fatalf("username = %q, err = %v", username, err)
	}
}

func TestWhoamiInvalidToken(t *testing.T) {
	srv, _ := newFakeAccountServer(t, "good-token", nil)

	if _, err := whoami(srv.URL, "wrong-token"); err == nil {
		t.Fatal("expected an error for an invalid token")
	}
}

func TestSaveAndLoadAccountLink(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".breaklab", "account.json")

	if link, err := loadAccountLink(path); err != nil || link != nil {
		t.Fatalf("loading a missing file: link = %+v, err = %v, want (nil, nil)", link, err)
	}

	want := &accountLink{Server: "https://example.com", Token: "abc123"}
	if err := saveAccountLink(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := loadAccountLink(path)
	if err != nil || got == nil || *got != *want {
		t.Fatalf("loaded = %+v, err = %v, want %+v", got, err, want)
	}
}

func TestRunLoginMissingFlagsExitsNonZeroWithoutStoring(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".breaklab", "account.json")
	withAccountLinkPath(t, path)

	if code := runLogin(nil); code != 2 {
		t.Fatalf("runLogin with no flags: code = %d, want 2", code)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("account link file should not exist, stat err = %v", err)
	}

	if code := runLogin([]string{"--server", "https://example.com"}); code != 2 {
		t.Fatalf("runLogin with only --server: code = %d, want 2", code)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("account link file should not exist, stat err = %v", err)
	}
}

func TestRunLoginSuccess(t *testing.T) {
	srv, _ := newFakeAccountServer(t, "good-token", nil)
	path := filepath.Join(t.TempDir(), ".breaklab", "account.json")
	withAccountLinkPath(t, path)

	if code := runLogin([]string{"--server", srv.URL, "--token", "good-token"}); code != 0 {
		t.Fatalf("runLogin: code = %d", code)
	}
	got, err := loadAccountLink(path)
	if err != nil || got == nil || got.Server != srv.URL || got.Token != "good-token" {
		t.Fatalf("stored link = %+v, err = %v", got, err)
	}
}

func TestLoginInvalidTokenStoresNothingAndLeavesExistingLinkUntouched(t *testing.T) {
	srv, _ := newFakeAccountServer(t, "good-token", nil)
	path := filepath.Join(t.TempDir(), ".breaklab", "account.json")

	existing := &accountLink{Server: srv.URL, Token: "good-token"}
	if err := saveAccountLink(path, existing); err != nil {
		t.Fatal(err)
	}

	if _, err := whoami(srv.URL, "wrong-token"); err == nil {
		t.Fatal("expected whoami to fail for the wrong token")
	}
	// runLogin only calls saveAccountLink after a successful whoami, so
	// the existing link must still be exactly what it was.
	got, err := loadAccountLink(path)
	if err != nil || got == nil || *got != *existing {
		t.Fatalf("existing link changed after a failed login attempt: got %+v, err %v", got, err)
	}
}
