package account

import (
	"encoding/json"
	"net/http"
	"strings"

	"breaklab/internal/account"
	"breaklab/internal/progress"
	"breaklab/internal/web/webcore"
)

// bearerUser resolves the request's `Authorization: Bearer <token>`
// header to its hosted account, for the /sync API — distinct from the
// cookie-based webcore.CurrentUser used by every other account-aware
// page.
func bearerUser(cfg webcore.Config, r *http.Request) (account.User, bool) {
	if cfg.AccountStore == nil {
		return account.User{}, false
	}
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || token == "" {
		return account.User{}, false
	}
	u, err := cfg.AccountStore.UserByToken(token)
	if err != nil {
		return account.User{}, false
	}
	return u, true
}

type syncAttemptsPayload struct {
	Attempts []progress.Attempt `json:"attempts"`
}

type whoamiPayload struct {
	Username string `json:"username"`
}

// syncWhoamiHandler resolves the bearer token to its account's
// username — how `breaklab login` confirms which account a token
// linked to, without echoing back anything more sensitive (email).
func syncWhoamiHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := bearerUser(cfg, r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(whoamiPayload{Username: u.Username})
	}
}

// syncPullHandler returns every attempt synced to the authenticated
// account.
func syncPullHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := bearerUser(cfg, r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		attempts, err := cfg.AccountStore.PullAttempts(u.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(syncAttemptsPayload{Attempts: attempts})
	}
}

// syncPushHandler stores every attempt in the request body against the
// authenticated account, ignoring any whose ID it already has.
func syncPushHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := bearerUser(cfg, r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var payload syncAttemptsPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := cfg.AccountStore.PushAttempts(u.ID, payload.Attempts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
