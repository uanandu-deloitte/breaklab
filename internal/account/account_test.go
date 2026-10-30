package account

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"breaklab/internal/progress"
)

func newTestStore(t *testing.T) Store {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "account.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestCreateUserSuccess(t *testing.T) {
	store := newTestStore(t)

	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == 0 || u.Email != "ada@example.com" || u.Username != "ada" {
		t.Fatalf("user = %+v", u)
	}
}

func TestCreateUserDuplicateEmailRejected(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.CreateUser("ada@example.com", "ada", "correct-password"); err != nil {
		t.Fatal(err)
	}

	_, err := store.CreateUser("ada@example.com", "ada2", "another-password")
	if !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("err = %v, want ErrEmailTaken", err)
	}
}

func TestCreateUserDuplicateUsernameRejected(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.CreateUser("ada@example.com", "ada", "correct-password"); err != nil {
		t.Fatal(err)
	}

	_, err := store.CreateUser("ada2@example.com", "ada", "another-password")
	if !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("err = %v, want ErrUsernameTaken", err)
	}
}

func TestCreateUserInvalidUsernameRejected(t *testing.T) {
	store := newTestStore(t)

	for _, bad := range []string{"Ada", "ada lovelace", "ada_lovelace", ""} {
		if _, err := store.CreateUser("ada@example.com", bad, "correct-password"); !errors.Is(err, ErrInvalidUsername) {
			t.Errorf("username %q: err = %v, want ErrInvalidUsername", bad, err)
		}
	}
}

func TestAuthenticateCorrectCredentials(t *testing.T) {
	store := newTestStore(t)
	created, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	u, err := store.Authenticate("ada@example.com", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != created.ID {
		t.Errorf("authenticated user ID = %d, want %d", u.ID, created.ID)
	}
}

func TestAuthenticateWrongPasswordRejected(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.CreateUser("ada@example.com", "ada", "correct-password"); err != nil {
		t.Fatal(err)
	}

	_, err := store.Authenticate("ada@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v, want ErrInvalidCredentials", err)
	}
}

func TestAuthenticateUnknownEmailRejectedSameAsWrongPassword(t *testing.T) {
	store := newTestStore(t)

	_, err := store.Authenticate("nobody@example.com", "whatever")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v, want ErrInvalidCredentials (same error as a wrong password, so the response can't reveal which email exists)", err)
	}
}

func TestPublicProfileDefaultsOffAndToggles(t *testing.T) {
	store := newTestStore(t)
	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if u.PublicProfile {
		t.Error("PublicProfile should default to false on signup")
	}
	got, err := store.Authenticate("ada@example.com", "correct-password")
	if err != nil || got.PublicProfile {
		t.Fatalf("freshly created account: PublicProfile = %v, err = %v, want false", got.PublicProfile, err)
	}

	if err := store.SetPublicProfile(u.ID, true); err != nil {
		t.Fatal(err)
	}
	got, err = store.UserByUsername("ada")
	if err != nil || !got.PublicProfile {
		t.Fatalf("after enabling: PublicProfile = %v, err = %v, want true", got.PublicProfile, err)
	}

	if err := store.SetPublicProfile(u.ID, false); err != nil {
		t.Fatal(err)
	}
	got, err = store.UserByUsername("ada")
	if err != nil || got.PublicProfile {
		t.Fatalf("after disabling: PublicProfile = %v, err = %v, want false", got.PublicProfile, err)
	}
}

func TestSessionLifecycle(t *testing.T) {
	store := newTestStore(t)
	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	token, expiresAt, err := store.CreateSession(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" || !expiresAt.After(time.Now()) {
		t.Fatalf("token = %q, expiresAt = %v", token, expiresAt)
	}

	got, err := store.UserBySession(token)
	if err != nil || got.ID != u.ID {
		t.Fatalf("UserBySession = %+v, err = %v", got, err)
	}

	if err := store.DeleteSession(token); err != nil {
		t.Fatal(err)
	}
	if _, err := store.UserBySession(token); !errors.Is(err, ErrNotFound) {
		t.Fatalf("after logout, err = %v, want ErrNotFound", err)
	}

	// Deleting an already-invalid session is a no-op, not an error.
	if err := store.DeleteSession(token); err != nil {
		t.Fatalf("deleting an already-deleted session should not error, got %v", err)
	}
}

func TestSessionExpires(t *testing.T) {
	store := newTestStore(t).(*sqliteStore)
	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := store.CreateSession(u.ID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := store.db.Exec(
		`UPDATE sessions SET expires_at = ? WHERE token_hash = ?`,
		time.Now().UTC().Add(-time.Hour), hashToken(token),
	); err != nil {
		t.Fatal(err)
	}

	if _, err := store.UserBySession(token); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expired session: err = %v, want ErrNotFound", err)
	}
}

func TestAPITokenLifecycle(t *testing.T) {
	store := newTestStore(t)
	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	token, err := store.IssueToken(u.ID)
	if err != nil || token == "" {
		t.Fatalf("token = %q, err = %v", token, err)
	}

	got, err := store.UserByToken(token)
	if err != nil || got.ID != u.ID {
		t.Fatalf("UserByToken = %+v, err = %v", got, err)
	}

	if err := store.RevokeToken(u.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.UserByToken(token); !errors.Is(err, ErrNotFound) {
		t.Fatalf("after revoke, err = %v, want ErrNotFound", err)
	}
}

func TestDeleteUserCascades(t *testing.T) {
	store := newTestStore(t).(*sqliteStore)
	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	sessionToken, _, err := store.CreateSession(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	apiToken, err := store.IssueToken(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(
		`INSERT INTO synced_attempts (user_id, id, lab, variant, started_at, failed_checks, hints_used, postmortem, status, hypothesis) VALUES (?, 'a1', 'web-down', '', ?, 0, 0, '', 'solved', '')`,
		u.ID, time.Now().UTC(),
	); err != nil {
		t.Fatal(err)
	}

	if err := store.DeleteUser(u.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := store.Authenticate("ada@example.com", "correct-password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("account should no longer exist, got err = %v", err)
	}
	if _, err := store.UserByUsername("ada"); !errors.Is(err, ErrNotFound) {
		t.Errorf("UserByUsername after deletion: err = %v, want ErrNotFound (so /u/ada would 404)", err)
	}
	if _, err := store.UserBySession(sessionToken); !errors.Is(err, ErrNotFound) {
		t.Errorf("session should be gone after account deletion, err = %v", err)
	}
	if _, err := store.UserByToken(apiToken); !errors.Is(err, ErrNotFound) {
		t.Errorf("API token should be gone after account deletion, err = %v", err)
	}
	var attemptCount int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM synced_attempts WHERE user_id = ?`, u.ID).Scan(&attemptCount); err != nil {
		t.Fatal(err)
	}
	if attemptCount != 0 {
		t.Errorf("synced_attempts rows for deleted user = %d, want 0", attemptCount)
	}
}

func TestPushAttemptsIsIdempotent(t *testing.T) {
	store := newTestStore(t)
	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	solved := time.Now().UTC()
	attempt := progress.Attempt{ID: "a1", Lab: "web-down", Status: progress.StatusSolved, StartedAt: solved.Add(-time.Minute), SolvedAt: &solved}

	if err := store.PushAttempts(u.ID, []progress.Attempt{attempt}); err != nil {
		t.Fatal(err)
	}
	if err := store.PushAttempts(u.ID, []progress.Attempt{attempt}); err != nil {
		t.Fatal(err)
	}

	got, err := store.PullAttempts(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "a1" {
		t.Fatalf("attempts = %+v, want exactly one (re-push must not duplicate)", got)
	}
}

func TestPullAttemptsIsolatedPerAccount(t *testing.T) {
	store := newTestStore(t)
	a, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	b, err := store.CreateUser("bob@example.com", "bob", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.PushAttempts(a.ID, []progress.Attempt{{ID: "a1", Lab: "web-down", Status: progress.StatusSolved}}); err != nil {
		t.Fatal(err)
	}

	got, err := store.PullAttempts(b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("account B's attempts = %+v, want none of account A's", got)
	}
}

func TestAPITokenRegenerateReplacesOldOne(t *testing.T) {
	store := newTestStore(t)
	u, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	oldToken, err := store.IssueToken(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	newToken, err := store.IssueToken(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if oldToken == newToken {
		t.Fatal("regenerating should produce a different token")
	}

	if _, err := store.UserByToken(oldToken); !errors.Is(err, ErrNotFound) {
		t.Fatalf("old token: err = %v, want ErrNotFound (replaced)", err)
	}
	if got, err := store.UserByToken(newToken); err != nil || got.ID != u.ID {
		t.Fatalf("new token: got = %+v, err = %v", got, err)
	}
}
