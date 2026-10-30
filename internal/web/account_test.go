package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"breaklab/internal/account"
	"breaklab/internal/web/webcore"
)

// newTestAccountStore opens a fresh SQLite-backed Store under a
// t.TempDir(), closed automatically at test cleanup.
func newTestAccountStore(t *testing.T) account.Store {
	t.Helper()
	store, err := account.Open(filepath.Join(t.TempDir(), "account.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func accountTestConfig(t *testing.T) Config {
	t.Helper()
	return Config{
		TemplatesDir: t.TempDir(),
		FrontendDir:  t.TempDir(),
		AccountStore: newTestAccountStore(t),
	}
}

func sessionCookieFrom(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == webcore.SessionCookieName {
			return c
		}
	}
	return nil
}

func postJSON(t *testing.T, mux http.Handler, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// tokenPayload mirrors the account package's unexported type of the
// same name.
type tokenPayload struct {
	Token string `json:"token"`
}

func TestAccountAPISignupLoginTokenLogoutFlow(t *testing.T) {
	cfg := accountTestConfig(t)
	mux := Routes(cfg)

	rec := postJSON(t, mux, "/api/signup", `{"email":"ada@example.com","username":"ada","password":"correct-password"}`, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("signup: code=%d body=%s", rec.Code, rec.Body)
	}
	cookie := sessionCookieFrom(rec)
	if cookie == nil {
		t.Fatal("signup should set a session cookie")
	}

	req := httptest.NewRequest("GET", "/api/account", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "ada@example.com") {
		t.Fatalf("account: code=%d body=%s", rec.Code, rec.Body)
	}

	rec = postJSON(t, mux, "/api/account/token", "", cookie)
	var tok tokenPayload
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &tok) != nil || tok.Token == "" {
		t.Fatalf("issue token: code=%d body=%s", rec.Code, rec.Body)
	}

	rec = postJSON(t, mux, "/api/account/token/revoke", "", cookie)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("revoke token: code=%d body=%s", rec.Code, rec.Body)
	}

	rec = postJSON(t, mux, "/api/logout", "", cookie)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("logout: code=%d body=%s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest("GET", "/api/account", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("account after logout: code=%d, want 401 (session was deleted)", rec.Code)
	}
}

func TestLoginAPIWrongPasswordShowsErrorNotWhichEmail(t *testing.T) {
	cfg := accountTestConfig(t)
	mux := Routes(cfg)

	postJSON(t, mux, "/api/signup", `{"email":"ada@example.com","username":"ada","password":"correct-password"}`, nil)

	rec := postJSON(t, mux, "/api/login", `{"email":"ada@example.com","password":"wrong"}`, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password: code=%d body=%s", rec.Code, rec.Body)
	}

	recUnknown := postJSON(t, mux, "/api/login", `{"email":"nobody@example.com","password":"wrong"}`, nil)
	if rec.Body.String() != recUnknown.Body.String() {
		t.Errorf("unknown-email and wrong-password errors differ (%q vs %q); must be identical so a response can't reveal which emails exist",
			rec.Body.String(), recUnknown.Body.String())
	}
}

func TestAccountAPIDeleteRequiresUsernameConfirmation(t *testing.T) {
	cfg := accountTestConfig(t)
	mux := Routes(cfg)

	rec := postJSON(t, mux, "/api/signup", `{"email":"ada@example.com","username":"ada","password":"correct-password"}`, nil)
	cookie := sessionCookieFrom(rec)

	rec = postJSON(t, mux, "/api/account/delete", `{"confirm_username":"not-ada"}`, cookie)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("wrong confirmation: code=%d body=%s, want 400 without deleting", rec.Code, rec.Body)
	}
	// Still logged in — the account must still exist.
	req := httptest.NewRequest("GET", "/api/account", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("account should still exist after a failed delete confirmation, code = %d", rec.Code)
	}

	rec = postJSON(t, mux, "/api/account/delete", `{"confirm_username":"ada"}`, cookie)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("correct confirmation: code=%d body=%s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest("GET", "/api/account", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("account after real deletion: code=%d, want 401", rec.Code)
	}
}

func TestAccountAPIPublicProfileToggle(t *testing.T) {
	cfg := accountTestConfig(t)
	mux := Routes(cfg)

	rec := postJSON(t, mux, "/api/signup", `{"email":"ada@example.com","username":"ada","password":"correct-password"}`, nil)
	cookie := sessionCookieFrom(rec)

	u, err := cfg.AccountStore.UserByUsername("ada")
	if err != nil || u.PublicProfile {
		t.Fatalf("fresh signup: PublicProfile = %v, err = %v, want false", u.PublicProfile, err)
	}

	rec = postJSON(t, mux, "/api/account/public-profile", `{"enabled":true}`, cookie)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("enabling: code=%d body=%s", rec.Code, rec.Body)
	}
	u, err = cfg.AccountStore.UserByUsername("ada")
	if err != nil || !u.PublicProfile {
		t.Fatalf("after enabling: PublicProfile = %v, err = %v, want true", u.PublicProfile, err)
	}

	rec = postJSON(t, mux, "/api/account/public-profile", `{"enabled":false}`, cookie)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("disabling: code=%d", rec.Code)
	}
	u, err = cfg.AccountStore.UserByUsername("ada")
	if err != nil || u.PublicProfile {
		t.Fatalf("after disabling: PublicProfile = %v, err = %v, want false", u.PublicProfile, err)
	}
}

func TestAccountRoutesNotMountedWithoutAccountStore(t *testing.T) {
	cfg := testConfig(t, 0)
	mux := Routes(cfg)

	for _, path := range []string{"/signup", "/login", "/account"} {
		req := httptest.NewRequest("GET", path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s with no AccountStore: code = %d, want 404", path, rec.Code)
		}
	}
	for _, path := range []string{"/api/signup", "/api/login", "/api/account"} {
		req := httptest.NewRequest("GET", path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s with no AccountStore: code = %d, want 404", path, rec.Code)
		}
	}
}
