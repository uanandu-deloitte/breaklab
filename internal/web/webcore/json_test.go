package webcore

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, map[string]int{"count": 3})

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var got map[string]int
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || got["count"] != 3 {
		t.Fatalf("body = %q, err = %v", rec.Body.String(), err)
	}
}

func TestWriteJSONError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSONError(rec, http.StatusNotFound, "not found")

	if rec.Code != http.StatusNotFound {
		t.Errorf("code = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var got apiErrorPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || got.Error != "not found" {
		t.Fatalf("body = %q, err = %v", rec.Body.String(), err)
	}
}
