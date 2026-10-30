package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"breaklab/internal/postmortem"
	"breaklab/internal/progress"
)

func solveAlpha(t *testing.T, r *Runner, ctx context.Context) {
	t.Helper()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestPostmortemCreatesExpectedFile(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)

	s, _ := LoadSession(r.SessionPath)
	pf, _ := progress.Load(r.ProgressPath)
	attempt := pf.Attempts[s.AttemptIndex]
	wantDate := attempt.SolvedAt.Local().Format("2006-01-02")
	wantPath := filepath.Join(r.PostmortemsDir, wantDate+"-alpha-first-variant.md")

	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, want := range []string{
		"lab: alpha", "variant: first-variant", "stage: linux",
		"date: " + wantDate, "hints_used: 0", "publish: false",
		"## Summary", "## Prevention",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("content missing %q:\n%s", want, content)
		}
	}

	pf, _ = progress.Load(r.ProgressPath)
	if pf.Attempts[s.AttemptIndex].Postmortem != path {
		t.Errorf("attempt.Postmortem = %q, want %q", pf.Attempts[s.AttemptIndex].Postmortem, path)
	}
}

func TestPostmortemCarriesHypothesisIntoFile(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	if err := r.Start(ctx, "alpha", false); err != nil {
		t.Fatal(err)
	}
	if err := r.SetHypothesis(ctx, "a stale cache, most likely"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}

	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "a stale cache, most likely") {
		t.Errorf("rendered postmortem is missing the saved hypothesis:\n%s", data)
	}
}

func TestPostmortemCallsOpenEditor(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)

	var gotPath string
	r.OpenEditor = func(_ context.Context, path string) error {
		gotPath = path
		return nil
	}

	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != path {
		t.Errorf("OpenEditor called with %q, want %q", gotPath, path)
	}
}

func TestPostmortemNothingSolvedErrors(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	if _, err := r.Postmortem(context.Background()); err == nil || err.Error() != "Solve a lab first" {
		t.Fatalf("err = %v, want %q", err, "Solve a lab first")
	}
}

func TestPostmortemSecondCallNothingLeftErrors(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)

	if _, err := r.Postmortem(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Postmortem(ctx); err == nil || err.Error() != "Solve a lab first" {
		t.Fatalf("second call: err = %v, want %q", err, "Solve a lab first")
	}
}

func TestPostmortemSecondCallPicksNextEligibleAttempt(t *testing.T) {
	r, _, _ := setup(t, "alpha", "beta")
	ctx := context.Background()
	solveAlpha(t, r, ctx)

	if err := r.Start(ctx, "beta", true); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}

	first, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	second, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatalf("both calls returned the same path: %q", first)
	}
}

const filledPostmortemBody = `## Summary
The homepage stopped responding after a routine maintenance window and
customers saw connection errors for about twenty minutes this morning.

## Detection
I ran curl against the site and got a connection refused error, then
checked whether nginx was running at all with systemctl status nginx.

## Root cause
A trailing semicolon was missing from the listen directive in the nginx
site config, so nginx refused to start after the maintenance script ran.

## Fix
I added the missing semicolon back to the listen directive, ran nginx
dash t to check the config, and then started nginx again successfully.

## Prevention
Add a CI step that runs nginx config test on every change to the nginx
configuration files before they are allowed to merge into the main branch.
`

func fillInSections(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw, _, err := postmortem.SplitFrontMatter(string(data))
	if err != nil {
		t.Fatal(err)
	}
	filled := "---\n" + raw + "\n---\n\n# Title\n\n" + filledPostmortemBody
	if err := os.WriteFile(path, []byte(filled), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestValidatePostmortemNoPathResolvesToLatestCreated(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)

	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fillInSections(t, path)

	res, err := r.ValidatePostmortem(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Valid || len(res.Failing) != 0 {
		t.Fatalf("res = %+v", res)
	}

	s, _ := LoadSession(r.SessionPath)
	pf, _ := progress.Load(r.ProgressPath)
	if pf.Attempts[s.AttemptIndex].Status != progress.StatusComplete {
		t.Errorf("attempt status = %q, want complete", pf.Attempts[s.AttemptIndex].Status)
	}
}

func TestValidatePostmortemInvalidDoesNotPromote(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)

	if _, err := r.Postmortem(ctx); err != nil {
		t.Fatal(err)
	}

	res, err := r.ValidatePostmortem(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if res.Valid || len(res.Failing) != 5 {
		t.Fatalf("res = %+v, want all 5 sections failing", res)
	}

	s, _ := LoadSession(r.SessionPath)
	pf, _ := progress.Load(r.ProgressPath)
	if pf.Attempts[s.AttemptIndex].Status != progress.StatusSolved {
		t.Errorf("attempt status = %q, want still solved", pf.Attempts[s.AttemptIndex].Status)
	}
}

func TestValidatePostmortemExplicitPathPromotesNonLatest(t *testing.T) {
	r, _, _ := setup(t, "alpha", "beta")
	ctx := context.Background()
	solveAlpha(t, r, ctx)
	alphaPath, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fillInSections(t, alphaPath)

	if err := r.Start(ctx, "beta", true); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Check(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Postmortem(ctx); err != nil {
		t.Fatal(err)
	}

	// alphaPath is no longer the latest-created postmortem (beta's is),
	// but an explicit path should still validate and promote it.
	res, err := r.ValidatePostmortem(ctx, alphaPath)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Valid {
		t.Fatalf("res = %+v", res)
	}

	pf, _ := progress.Load(r.ProgressPath)
	var found bool
	for _, a := range pf.Attempts {
		if a.Postmortem == alphaPath {
			found = true
			if a.Status != progress.StatusComplete {
				t.Errorf("alpha's attempt status = %q, want complete", a.Status)
			}
		}
	}
	if !found {
		t.Fatal("alpha's attempt not found")
	}
}

func TestValidatePostmortemAlreadyCompleteIsNoOp(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)
	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fillInSections(t, path)

	if _, err := r.ValidatePostmortem(ctx, ""); err != nil {
		t.Fatal(err)
	}
	res, err := r.ValidatePostmortem(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Valid {
		t.Fatalf("re-validating an already-complete postmortem: res = %+v", res)
	}
}

func TestValidatePostmortemNoneToValidateErrors(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	if _, err := r.ValidatePostmortem(context.Background(), ""); err == nil {
		t.Fatal("expected an error when no postmortem has ever been created")
	}
}

// filledPostmortemBodyMatchingVariant is filledPostmortemBody with its
// Root cause section rewritten to actually overlap with "alpha"'s
// "first-variant" fixture hint ("first-variant hint" — see
// runner_test.go's labYAML), instead of the unrelated nginx story.
const filledPostmortemBodyMatchingVariant = `## Summary
The homepage stopped responding after a routine maintenance window and
customers saw connection errors for about twenty minutes this morning.

## Detection
I ran curl against the site and got a connection refused error, then
checked whether nginx was running at all with systemctl status nginx.

## Root cause
Following the first-variant hint pointed straight at the actual cause
of this incident, which was exactly what that hint described happening.

## Fix
I added the missing semicolon back to the listen directive, ran nginx
dash t to check the config, and then started nginx again successfully.

## Prevention
Add a CI step that runs nginx config test on every change to the nginx
configuration files before they are allowed to merge into the main branch.
`

func fillInSectionsWithBody(t *testing.T, path, body string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw, _, err := postmortem.SplitFrontMatter(string(data))
	if err != nil {
		t.Fatal(err)
	}
	filled := "---\n" + raw + "\n---\n\n# Title\n\n" + body
	if err := os.WriteFile(path, []byte(filled), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestValidatePostmortemRootCauseMismatchWarnsButStillCompletes(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)
	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// filledPostmortemBody's root cause (a missing nginx semicolon) has
	// no overlap with "first-variant"'s fixture hint.
	fillInSections(t, path)

	res, err := r.ValidatePostmortem(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Valid {
		t.Fatalf("res = %+v, want valid: the warning is advisory, never blocking", res)
	}
	if res.RootCauseWarning == "" || !strings.Contains(res.RootCauseWarning, "variant") {
		t.Errorf("RootCauseWarning = %q, want it to name the expected hint words", res.RootCauseWarning)
	}
}

func TestValidatePostmortemRootCauseMatchHasNoWarning(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)
	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fillInSectionsWithBody(t, path, filledPostmortemBodyMatchingVariant)

	res, err := r.ValidatePostmortem(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Valid {
		t.Fatalf("res = %+v, want valid", res)
	}
	if res.RootCauseWarning != "" {
		t.Errorf("RootCauseWarning = %q, want empty: root cause overlaps the variant's hint", res.RootCauseWarning)
	}
}

func TestPostmortemFilenameCollisionGetsSuffix(t *testing.T) {
	r, _, _ := setup(t, "alpha")
	ctx := context.Background()
	solveAlpha(t, r, ctx)

	s, _ := LoadSession(r.SessionPath)
	pf, _ := progress.Load(r.ProgressPath)
	date := pf.Attempts[s.AttemptIndex].SolvedAt.Local().Format("2006-01-02")

	collidingPath := filepath.Join(r.PostmortemsDir, date+"-alpha-first-variant.md")
	if err := os.MkdirAll(r.PostmortemsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(collidingPath, []byte("already here\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	path, err := r.Postmortem(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(r.PostmortemsDir, date+"-alpha-first-variant-2.md")
	if path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
}
