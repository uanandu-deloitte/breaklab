# Solution: Read the Pipeline

```bash
cat <<'EOF' > /workspace/answers.env
TRIGGER_EVENTS=push,workflow_dispatch
PARALLEL_JOBS=lint,build
WAITING_JOB=test
WAITING_ON=build
CACHE_KEY_FILE=requirements.lock
MASKED_FAILURE_JOB=test
SWALLOWED_STEP="Run linter"
EOF
```

1. `on:` lists two triggers: `push` (restricted to `branches: [main]`) and
   `workflow_dispatch` (a manual run). Either one starts this workflow.

2. `lint` and `build` both have no `needs:` key, so neither waits on the
   other — they start together and run in parallel. `test` has
   `needs: build`, naming exactly one other job, so it waits for `build`
   to finish. (`deploy` also waits, but on two jobs — `needs: [build,
   test]` — not exactly one, which is why it isn't the answer to this
   question.)

3. The `build` job's cache step keys on
   `hashFiles('requirements.lock')`. If `requirements.lock` changes, its
   hash changes, so the computed key changes too — the next run can't
   find a cache entry under the new key, which is exactly the point:
   changed dependencies should invalidate the old cache rather than
   reuse stale installed packages.

4. `test` is the job marked `continue-on-error: true`. Its own failure
   gets reported, but nothing downstream treats it as a workflow
   failure — `deploy` still runs even if `test` failed, because
   `continue-on-error` masks that failure from everything that depends
   on it.

5. The `lint` job's `Run linter` step runs `./scripts/lint.sh || true`.
   Whatever `lint.sh` actually exits with, `|| true` replaces it with
   `0`, so that step always reports success regardless of whether the
   linter found anything.
