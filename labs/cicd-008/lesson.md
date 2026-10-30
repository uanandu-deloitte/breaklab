# A job output only exists if a step actually sets it

Declaring `outputs: { ready: ${{ steps.build.outputs.ready }} }` on a
job doesn't create that output — it just wires up where to find it if a
step produces it. If no step ever writes to `$GITHUB_OUTPUT`, the
output is permanently empty, and anything that reads it downstream
(`needs.build.outputs.ready`) gets an empty string back, silently.

## Why this matters

`if: ${{ needs.build.outputs.ready == 'true' }}` with an empty string on
the left is just `'' == 'true'` — always false. The job is skipped, not
failed — no red X, no error, nothing in the logs to point at. This is
one of the most common real GitHub Actions bugs: a `needs.<job>.outputs.*`
reference that looks wired up but was never actually connected to a step
that sets it.

## Where to look

- `.github/workflows/ci.yml` — read `build`'s steps for anywhere that
  should write to `$GITHUB_OUTPUT` and doesn't.
- `act push -P ubuntu-latest=-self-hosted` — runs the real workflow
  locally; watch which jobs it reports running at all.

## The fix

Set the output the job already declares, from the step that actually
knows the real answer:

```yaml
- name: Build the app
  id: build
  run: |
    mkdir -p /workspace/dist
    echo "built" > /workspace/dist/app.txt
    echo "ready=true" >> "$GITHUB_OUTPUT"
```
