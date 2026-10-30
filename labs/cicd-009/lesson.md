# `rules:` decides whether a job exists in the pipeline at all

A GitLab CI job's `rules:` list is evaluated before the job ever runs.
If none of its `if:` conditions are true, the job isn't "failed" — it's
never even created for that pipeline. No red X, no log, nothing to
click on. The only way to notice is that the job is simply missing from
the pipeline view.

## Why this matters

`$CI_COMMIT_TAG` is one of GitLab's predefined variables, but it's only
ever set on a pipeline triggered by pushing a **tag** — on an ordinary
branch push it's always empty. A `rules: - if: '$CI_COMMIT_TAG'` guard
on a job you expect to run on every push is a silent trap: it works
exactly once, the first time someone tests it with a tag, and then
quietly does nothing for every normal push afterward.

## Where to look

- `.gitlab-ci.yml` — read `deploy`'s `rules:` and ask: what variable is
  this actually checking, and is it ever set for the trigger type this
  pipeline runs under?
- `gitlab-ci-local --force-shell-executor` — runs the real pipeline
  locally; watch which jobs it reports at all.

## The fix

Match the condition to the trigger you actually care about —
`$CI_COMMIT_BRANCH` is set on every branch push:

```yaml
deploy:
  rules:
    - if: '$CI_COMMIT_BRANCH'
```
