package web

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"breaklab/internal/progress"
)

// TestPublicProfileAPINeverLeaksPostmortemPathOrDraftContent confirms
// that an attempt's Postmortem field — a local file path the hosted
// server never reads the content or publish-flag of — never appears
// anywhere in GET /api/u/{name}'s response. No postmortem content is
// ever synced to the hosted side at all (design.md), so the only way
// this leaks is the raw path string itself ending up in the response;
// this pins that it does not.
func TestPublicProfileAPINeverLeaksPostmortemPathOrDraftContent(t *testing.T) {
	cfg := publicProfileTestConfig(t)

	u, err := cfg.AccountStore.CreateUser("ada@example.com", "ada", "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.AccountStore.SetPublicProfile(u.ID, true); err != nil {
		t.Fatal(err)
	}

	const draftMarker = "DRAFT_POSTMORTEM_PATH_MARKER"
	solved := time.Now().UTC()
	if err := cfg.AccountStore.PushAttempts(u.ID, []progress.Attempt{
		{
			ID: "h1", Lab: "alpha", Status: progress.StatusComplete,
			StartedAt: solved.Add(-time.Minute), SolvedAt: &solved,
			Postmortem: "postmortems/" + draftMarker + ".md",
			Hypothesis: "a hypothesis that must never be public",
		},
	}); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	Routes(cfg).ServeHTTP(rec, httptest.NewRequest("GET", "/api/u/ada", nil))
	if rec.Code != 200 {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	if strings.Contains(body, draftMarker) {
		t.Error("the postmortem's file path leaked into the public profile API response")
	}
	if strings.Contains(body, "a hypothesis that must never be public") {
		t.Error("the attempt's private Hypothesis text leaked into the public profile API response")
	}
}

// The old TestPublicProfileTemplateNeverReferencesVariantOrPostmortemFields
// statically checked that public_profile.html's template source never
// referenced a field that could carry a variant id, solution script,
// hint text, or postmortem path/href. That template is gone (migrated
// to the `frontend` capability); the property it pinned is now a type
// guarantee instead of a template-source grep: publicProfileView and
// the site.StageView/site.Totals types it wraps have no Variant,
// Hypothesis, or Hint field at all, so no JSON response from this
// package can ever carry one, and PostmortemHref is already covered by
// the leak test above (always empty — BuildIndex is called with an
// empty "published" map).
