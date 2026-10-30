---
lab: {{.Lab}}
variant: {{.Variant}}
stage: {{.Stage}}
date: {{.Date}}
minutes_to_fix: {{.MinutesToFix}}
hints_used: {{.HintsUsed}}
publish: false
---

# {{.Title}}

## Hypothesis
{{if .Hypothesis}}{{.Hypothesis}}{{else}}<!-- What did you think was broken before you started digging? Optional, but worth comparing against Root cause below. -->{{end}}

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
