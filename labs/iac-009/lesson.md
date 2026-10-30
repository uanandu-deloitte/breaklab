# Plan diffs config against state — not against the file on disk

Terraform's `plan` command answers one question: does the **declared
config** match what Terraform's own **state** last recorded? It doesn't
matter what the real file on disk currently says — if state and config
disagree, `plan` reports a pending change, and if you hand-edit the real
file to "fix" it without ever running `apply`, state still disagrees
with config, and the very next `plan` will say so.

## Why this matters

This is one of the most common real incidents on any team using
Terraform: someone merges a change to a `.tf` file, the PR is approved,
but the pipeline (or person) that's supposed to run `terraform apply`
for it never does. The real world keeps running on the *old* applied
state until someone actually applies the new config — Terraform never
applies anything just because a file changed on disk.

## Where to look

- `/infra/main.tf` — the currently declared config.
- `/etc/motd-app/banner.txt` — the real file `local_file.banner`
  manages.
- `terraform plan` (run from `/infra`) — the authoritative answer to
  "what, if anything, is out of sync."

## The fix

```
cd /infra
terraform plan                       # see exactly what's pending, and why
terraform apply -auto-approve        # reconcile the real world to match main.tf
terraform plan -detailed-exitcode    # exit 0 once nothing's left pending
```

Never hand-edit a file a `resource` block manages directly — the next
`apply` (by you, a teammate, or a pipeline) will just overwrite it back
to whatever the state file remembers, silently undoing your fix.
