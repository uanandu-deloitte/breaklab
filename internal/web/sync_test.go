package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"breaklab/internal/account"
	"breaklab/internal/progress"
)

// syncAttemptsPayload mirrors the account package's unexported type of
// the same name.
type syncAttemptsPayload struct {
	Attempts []progress.Attempt `json:"attempts"`
}

func createAccountAndToken(t *testing.T, store account.Store, email, username string) string {
	t.Helper()
	u, err := store.CreateUser(email, username, "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	token, err := store.IssueToken(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func pushAttempts(t *testing.T, mux http.Handler, token string, attempts []progress.Attempt) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(syncAttemptsPayload{Attempts: attempts})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/sync/attempts", bytes.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func pullAttempts(t *testing.T, mux http.Handler, token string) (*httptest.ResponseRecorder, syncAttemptsPayload) {
	t.Helper()
	req := httptest.NewRequest("GET", "/sync/attempts", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var payload syncAttemptsPayload
	if rec.Code == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
	}
	return rec, payload
}

func TestSyncRejectsMissingOrInvalidToken(t *testing.T) {
	cfg := accountTestConfig(t)
	mux := Routes(cfg)

	if rec, _ := pullAttempts(t, mux, ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("no token: code = %d, want 401", rec.Code)
	}
	if rec, _ := pullAttempts(t, mux, "not-a-real-token"); rec.Code != http.StatusUnauthorized {
		t.Errorf("invalid token: code = %d, want 401", rec.Code)
	}
	if rec := pushAttempts(t, mux, "not-a-real-token", nil); rec.Code != http.StatusUnauthorized {
		t.Errorf("push with invalid token: code = %d, want 401", rec.Code)
	}
}

func TestSyncCrossAccountIsolation(t *testing.T) {
	cfg := accountTestConfig(t)
	mux := Routes(cfg)

	tokenA := createAccountAndToken(t, cfg.AccountStore, "ada@example.com", "ada")
	tokenB := createAccountAndToken(t, cfg.AccountStore, "bob@example.com", "bob")

	solved := time.Now().UTC()
	rec := pushAttempts(t, mux, tokenA, []progress.Attempt{
		{ID: "a1", Lab: "web-down", Status: progress.StatusSolved, StartedAt: solved.Add(-time.Minute), SolvedAt: &solved},
	})
	if rec.Code != http.StatusNoContent {
		t.Fatalf("push: code = %d, body = %s", rec.Code, rec.Body)
	}

	_, payloadB := pullAttempts(t, mux, tokenB)
	if len(payloadB.Attempts) != 0 {
		t.Errorf("account B pulled %+v, want none of account A's attempts", payloadB.Attempts)
	}

	_, payloadA := pullAttempts(t, mux, tokenA)
	if len(payloadA.Attempts) != 1 || payloadA.Attempts[0].ID != "a1" {
		t.Errorf("account A pulled %+v, want its own attempt a1", payloadA.Attempts)
	}
}

func TestSyncPushIsIdempotent(t *testing.T) {
	cfg := accountTestConfig(t)
	mux := Routes(cfg)
	token := createAccountAndToken(t, cfg.AccountStore, "ada@example.com", "ada")

	solved := time.Now().UTC()
	attempts := []progress.Attempt{{ID: "a1", Lab: "web-down", Status: progress.StatusSolved, StartedAt: solved.Add(-time.Minute), SolvedAt: &solved}}

	for i := 0; i < 2; i++ {
		if rec := pushAttempts(t, mux, token, attempts); rec.Code != http.StatusNoContent {
			t.Fatalf("push %d: code = %d, body = %s", i, rec.Code, rec.Body)
		}
	}

	_, payload := pullAttempts(t, mux, token)
	if len(payload.Attempts) != 1 {
		t.Fatalf("attempts after re-push = %+v, want exactly one", payload.Attempts)
	}
}
