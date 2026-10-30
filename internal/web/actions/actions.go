package actions

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
)

// ErrActionsUnavailable means the server was configured without a Docker
// client (e.g. most test Configs), so the lab action routes can't run.
var ErrActionsUnavailable = errors.New("lab actions are not available: no Docker client configured")

// newActionRunner builds a Runner over cfg's own paths — never the
// process's default repo-relative paths — so action routes stay
// confined to whatever a test Config points at, exactly like every
// read-only view in this package. Its Out is discarded: the CLI-style
// text Start/Stop write there duplicates what the lab detail view
// already renders from disk.
func newActionRunner(cfg webcore.Config) (*runner.Runner, error) {
	if cfg.Docker == nil {
		return nil, ErrActionsUnavailable
	}
	r := runner.New(cfg.Docker)
	r.LabsDir = cfg.LabsDir
	r.RoadmapPath = cfg.RoadmapPath
	r.SessionPath = cfg.SessionPath
	r.ProgressPath = cfg.ProgressPath
	if cfg.HistoryPath != "" {
		r.HistoryPath = cfg.HistoryPath
	}
	r.Out = io.Discard
	return r, nil
}

// actionErrorStatus maps a runner action error to a response status:
// ErrActionsUnavailable (no Docker configured) is a 503, everything
// else (a blocked start, a failed check, ...) is a plain 500 — the
// Astro page shows the message either way, it doesn't branch on code.
func actionErrorStatus(err error) int {
	if errors.Is(err, ErrActionsUnavailable) {
		return http.StatusServiceUnavailable
	}
	return http.StatusInternalServerError
}

// labStartAPIHandler launches labID's container: the only route that
// starts Docker. A plain click never displaces a different lab's
// already-running container — Runner.Start returns an error naming it
// instead, and the client offers an explicit confirm (body {"confirm":
// true}) that passes through as --yes does on the CLI.
func labStartAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := webcore.RequireActionUser(cfg, w, r); !ok {
			return
		}
		id := r.PathValue("id")
		rn, err := newActionRunner(cfg)
		if err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		var body struct {
			Confirm bool `json:"confirm"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body) // no/invalid body just means confirm=false
		if err := rn.Start(r.Context(), id, body.Confirm); err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// labCheckAPIHandler runs the active lab's check and returns the full
// runner.CheckResult (PASS/FAIL message, and per-check results for a
// multi-check lab) directly as this response's body.
func labCheckAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := webcore.RequireActionUser(cfg, w, r); !ok {
			return
		}
		rn, err := newActionRunner(cfg)
		if err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		res, err := rn.Check(r.Context())
		if err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		webcore.WriteJSON(w, res)
	}
}

// hintPayload is labHintAPIHandler's response. The hint text itself is
// never persisted (design.md Decision 4) — it only ever exists for
// this one response.
type hintPayload struct {
	Hint string `json:"hint"`
}

// labHintAPIHandler reveals the active lab's next unseen hint.
func labHintAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := webcore.RequireActionUser(cfg, w, r); !ok {
			return
		}
		rn, err := newActionRunner(cfg)
		if err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		hint, err := rn.Hint(r.Context())
		if err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		webcore.WriteJSON(w, hintPayload{Hint: hint})
	}
}

// labHypothesisAPIHandler records the learner's own guess at the cause
// on the active attempt (design.md/loop's "write it down before you
// run the check" step), from this request's JSON body.
func labHypothesisAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := webcore.RequireActionUser(cfg, w, r); !ok {
			return
		}
		rn, err := newActionRunner(cfg)
		if err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		var body struct {
			Hypothesis string `json:"hypothesis"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			webcore.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if err := rn.SetHypothesis(r.Context(), body.Hypothesis); err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// labStopAPIHandler removes the active lab's container and clears the
// session, freeing the next Start (on any lab) to run without a
// confirm step.
func labStopAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := webcore.RequireActionUser(cfg, w, r); !ok {
			return
		}
		rn, err := newActionRunner(cfg)
		if err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		if err := rn.Stop(r.Context()); err != nil {
			webcore.WriteJSONError(w, actionErrorStatus(err), err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// RegisterRoutes mounts the lab action JSON API routes on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("POST /api/labs/{id}/start", labStartAPIHandler(cfg))
	mux.HandleFunc("POST /api/labs/{id}/check", labCheckAPIHandler(cfg))
	mux.HandleFunc("POST /api/labs/{id}/hint", labHintAPIHandler(cfg))
	mux.HandleFunc("POST /api/labs/{id}/hypothesis", labHypothesisAPIHandler(cfg))
	mux.HandleFunc("POST /api/labs/{id}/stop", labStopAPIHandler(cfg))
}
