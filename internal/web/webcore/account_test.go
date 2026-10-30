package webcore

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"breaklab/internal/account"
)

func newTestAccountStore(t *testing.T) account.Store {
	t.Helper()
	store, err := account.Open(filepath.Join(t.TempDir(), "account.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestRequireActionUserLocalModeNoCookie(t *testing.T) {
	cfg := Config{}
	req := httptest.NewRequest("POST", "/api/labs/web-down/start", nil)
	rec := httptest.NewRecorder()

	u, ok := RequireActionUser(cfg, rec, req)
	if !ok {
		t.Fatalf("ok = false in local mode, want true (no AccountStore means every request is allowed)")
	}
	if u != (account.User{}) {
		t.Errorf("User = %+v, want zero value in local mode", u)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("response body = %q, want nothing written in the ok=true path", rec.Body.String())
	}
}

func TestRequireActionUserHostedNoCookie(t *testing.T) {
	cfg := Config{AccountStore: newTestAccountStore(t)}
	req := httptest.NewRequest("POST", "/api/labs/web-down/start", nil)
	rec := httptest.NewRecorder()

	_, ok := RequireActionUser(cfg, rec, req)
	if ok {
		t.Fatal("ok = true with no session in hosted mode, want false")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestRequireActionUserHostedValidCookie(t *testing.T) {
	store := newTestAccountStore(t)
	cfg := Config{AccountStore: store}
	created, err := store.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := store.CreateSession(created.ID)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/labs/web-down/start", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: token})
	rec := httptest.NewRecorder()

	u, ok := RequireActionUser(cfg, rec, req)
	if !ok {
		t.Fatalf("ok = false with a valid session, want true; body = %s", rec.Body)
	}
	if u.ID != created.ID {
		t.Errorf("resolved user ID = %d, want %d", u.ID, created.ID)
	}
}

func TestRequireSessionAPIAlwaysRequiresLogin(t *testing.T) {
	cfg := Config{AccountStore: newTestAccountStore(t)}
	req := httptest.NewRequest("GET", "/api/account", nil)
	rec := httptest.NewRecorder()

	_, ok := RequireSessionAPI(cfg, rec, req)
	if ok {
		t.Fatal("ok = true with no session, want false")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}
