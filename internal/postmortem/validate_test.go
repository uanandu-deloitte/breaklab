package postmortem

import (
	"reflect"
	"sort"
	"testing"
)

const filledBody = `
# The website returns nothing

## Summary
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

func TestValidateSectionsAllBlankFailAll(t *testing.T) {
	blank := `
# Title

## Summary
<!-- What broke and what was the impact, in two or three sentences. -->

## Detection
<!-- Which commands and outputs led you to the cause, in order. -->

## Root cause
<!-- The actual cause, not the symptom. -->

## Fix
<!-- What you changed, with the exact commands. -->

## Prevention
<!-- What would stop this in production: a test, an alert, a CI check. -->
`
	got := ValidateSections(blank)
	sort.Strings(got)
	want := append([]string{}, RequiredSections...)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want all five: %v", got, want)
	}
}

func TestValidateSectionsFilledPasses(t *testing.T) {
	if got := ValidateSections(filledBody); len(got) != 0 {
		t.Errorf("got failing = %v, want none", got)
	}
}

func TestValidateSectionsMissingHeadingFails(t *testing.T) {
	// "## Root cause" heading removed entirely, its text folded into Detection.
	body := `
# Title

## Summary
The homepage stopped responding after a routine maintenance window and
customers saw connection errors for about twenty minutes this morning.

## Detection
I ran curl against the site and got a connection refused error, then
checked whether nginx was running at all with systemctl status nginx.
A trailing semicolon was missing from the listen directive.

## Fix
I added the missing semicolon back to the listen directive, ran nginx
dash t to check the config, and then started nginx again successfully.

## Prevention
Add a CI step that runs nginx config test on every change to the nginx
configuration files before they are allowed to merge into the main branch.
`
	got := ValidateSections(body)
	if len(got) != 1 || got[0] != "Root cause" {
		t.Errorf("got %v, want [Root cause]", got)
	}
}
