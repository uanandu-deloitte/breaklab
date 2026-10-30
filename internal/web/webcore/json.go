package webcore

import (
	"encoding/json"
	"net/http"
)

// apiErrorPayload is the one error shape every `frontend`-capability
// JSON endpoint uses, so client-side code has exactly one error
// format to handle regardless of which endpoint failed.
type apiErrorPayload struct {
	Error string `json:"error"`
}

// WriteJSON encodes v as the response body with the correct content
// type. Used for every successful `frontend`-capability API response.
func WriteJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// WriteJSONError writes status and msg in the shared apiErrorPayload
// shape. Used for every failed `frontend`-capability API response —
// never a plain-text http.Error for these endpoints, so client-side
// error handling never has to branch on content type.
func WriteJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiErrorPayload{Error: msg})
}
