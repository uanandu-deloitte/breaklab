package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"breaklab/internal/account"
	"breaklab/internal/web/webcore"
)

// requireSessionAPI resolves the current user or writes a 401 JSON
// error; callers must return immediately when ok is false.
func requireSessionAPI(cfg webcore.Config, w http.ResponseWriter, r *http.Request) (account.User, bool) {
	u, ok := webcore.CurrentUser(cfg, r)
	if !ok {
		webcore.WriteJSONError(w, http.StatusUnauthorized, "not logged in")
	}
	return u, ok
}

// signupErrMessage maps a Store error to the text shown on the signup
// form; anything else (a DB error) gets a generic message rather than
// leaking internals.
func signupErrMessage(err error) string {
	switch {
	case errors.Is(err, account.ErrEmailTaken):
		return "That email is already registered."
	case errors.Is(err, account.ErrUsernameTaken):
		return "That username is already taken."
	case errors.Is(err, account.ErrInvalidUsername):
		return "Usernames can only contain lowercase letters, digits, and hyphens."
	default:
		return "Could not create account."
	}
}

type credentialsPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// signupAPIHandler creates an account and sets the session cookie, the
// same as the old form-based handler, with a JSON error instead of a
// redirect+query-flash on failure.
func signupAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body credentialsPayload
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			webcore.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		u, err := cfg.AccountStore.CreateUser(body.Email, body.Username, body.Password)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusBadRequest, signupErrMessage(err))
			return
		}
		token, expiresAt, err := cfg.AccountStore.CreateSession(u.ID)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		webcore.SetSessionCookie(w, r, token, expiresAt)
		webcore.WriteJSON(w, u)
	}
}

// loginAPIHandler authenticates and sets the session cookie. The error
// message never reveals whether the email itself is registered.
func loginAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body credentialsPayload
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			webcore.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		u, err := cfg.AccountStore.Authenticate(body.Email, body.Password)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusUnauthorized, "Incorrect email or password.")
			return
		}
		token, expiresAt, err := cfg.AccountStore.CreateSession(u.ID)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		webcore.SetSessionCookie(w, r, token, expiresAt)
		webcore.WriteJSON(w, u)
	}
}

// logoutAPIHandler clears the current session, if any.
func logoutAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(webcore.SessionCookieName); err == nil {
			_ = cfg.AccountStore.DeleteSession(c.Value)
		}
		webcore.ClearSessionCookie(w, r)
		w.WriteHeader(http.StatusNoContent)
	}
}

// accountAPIHandler returns the signed-in user's own account data.
func accountAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := requireSessionAPI(cfg, w, r)
		if !ok {
			return
		}
		webcore.WriteJSON(w, u)
	}
}

type tokenPayload struct {
	Token string `json:"token"`
}

// accountIssueTokenAPIHandler issues a new API token, returned directly
// in this response's body — the one time its plaintext value exists
// (design.md Decision 4's "shown once" pattern, same as a lab hint).
func accountIssueTokenAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := requireSessionAPI(cfg, w, r)
		if !ok {
			return
		}
		token, err := cfg.AccountStore.IssueToken(u.ID)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		webcore.WriteJSON(w, tokenPayload{Token: token})
	}
}

func accountRevokeTokenAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := requireSessionAPI(cfg, w, r)
		if !ok {
			return
		}
		if err := cfg.AccountStore.RevokeToken(u.ID); err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func accountPublicProfileAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := requireSessionAPI(cfg, w, r)
		if !ok {
			return
		}
		var body struct {
			Enabled bool `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			webcore.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if err := cfg.AccountStore.SetPublicProfile(u.ID, body.Enabled); err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// accountDeleteAPIHandler deletes the account once confirm_username
// matches exactly, same as the old handler's confirmation rule.
func accountDeleteAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := requireSessionAPI(cfg, w, r)
		if !ok {
			return
		}
		var body struct {
			ConfirmUsername string `json:"confirm_username"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			webcore.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if body.ConfirmUsername != u.Username {
			webcore.WriteJSONError(w, http.StatusBadRequest, "Type your username exactly to confirm deletion.")
			return
		}
		if err := cfg.AccountStore.DeleteUser(u.ID); err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		webcore.ClearSessionCookie(w, r)
		w.WriteHeader(http.StatusNoContent)
	}
}

// RegisterRoutes mounts every hosted-account JSON API route (signup,
// login, logout, account settings, public profile, sync) on mux. Only
// called when cfg.AccountStore != nil.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("POST /api/signup", signupAPIHandler(cfg))
	mux.HandleFunc("POST /api/login", loginAPIHandler(cfg))
	mux.HandleFunc("POST /api/logout", logoutAPIHandler(cfg))
	mux.HandleFunc("GET /api/account", accountAPIHandler(cfg))
	mux.HandleFunc("POST /api/account/token", accountIssueTokenAPIHandler(cfg))
	mux.HandleFunc("POST /api/account/token/revoke", accountRevokeTokenAPIHandler(cfg))
	mux.HandleFunc("POST /api/account/public-profile", accountPublicProfileAPIHandler(cfg))
	mux.HandleFunc("POST /api/account/delete", accountDeleteAPIHandler(cfg))
	mux.HandleFunc("GET /api/u/{name}", publicProfileAPIHandler(cfg))
	mux.HandleFunc("GET /sync/attempts", syncPullHandler(cfg))
	mux.HandleFunc("POST /sync/attempts", syncPushHandler(cfg))
	mux.HandleFunc("GET /sync/whoami", syncWhoamiHandler(cfg))
}
