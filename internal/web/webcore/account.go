package webcore

import (
	"net/http"
	"time"

	"breaklab/internal/account"
	"breaklab/internal/progress"
)

// SessionCookieName is the hosted-account session cookie. It is only
// ever set/read when cfg.AccountStore is non-nil — the local,
// account-free server never touches cookies at all.
const SessionCookieName = "breaklab_session"

// SetSessionCookie sets the hosted-account session cookie.
func SetSessionCookie(w http.ResponseWriter, r *http.Request, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie clears the hosted-account session cookie.
func ClearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

// CurrentUser resolves the request's session cookie to its account.
// False whenever AccountStore is nil (local-only mode) or the cookie is
// missing, unknown, or expired.
func CurrentUser(cfg Config, r *http.Request) (account.User, bool) {
	if cfg.AccountStore == nil {
		return account.User{}, false
	}
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return account.User{}, false
	}
	u, err := cfg.AccountStore.UserBySession(c.Value)
	if err != nil {
		return account.User{}, false
	}
	return u, true
}

// ProgressForRequest returns pf merged with the current request's
// hosted account's synced attempts, if any (web-ui spec's
// account-aware dashboard/roadmap scenarios). With no session, no
// AccountStore, or a failed hosted lookup, it degrades to pf exactly
// as-is — the same behavior every handler had before accounts
// existed, never a 500.
func ProgressForRequest(cfg Config, r *http.Request, pf *progress.ProgressFile) *progress.ProgressFile {
	u, ok := CurrentUser(cfg, r)
	if !ok {
		return pf
	}
	hosted, err := cfg.AccountStore.PullAttempts(u.ID)
	if err != nil {
		return pf
	}
	return progress.MergeSynced(pf, hosted)
}
