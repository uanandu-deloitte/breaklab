# Read the Pipeline

There's no live pipeline run to watch here. Instead, `/workspace/pipeline.yml`
holds a GitHub-Actions-style workflow for a fictional project. It is never
executed — just read it, and answer the questions below from the text
alone.

## Questions

1. **Triggers.** Which event(s) under `on:` cause this workflow to run?
2. **Parallel vs. waiting.** Two jobs have no `needs:` at all, so they run
   in parallel with each other. A third job's `needs:` names exactly one
   other job, so it must wait for that one job to finish first. Which two
   jobs run in parallel, and which job waits — and on which job does it
   wait?
3. **Cache key.** The `build` job restores a cache keyed in part by
   `hashFiles(...)` on one file. Which file is that, and what happens to
   the cache key if that file's contents change?
4. **Masked failure, job level.** One job is marked
   `continue-on-error: true`. Which job?
5. **Swallowed failure, step level.** One step's `run:` command ends in
   `|| true`, which turns a real non-zero exit into `0` no matter what
   happened. Which step (by its exact `name:`) is that?

## Record your answers

Write these to `/workspace/answers.env`, one `KEY=value` per line:

```
TRIGGER_EVENTS=
PARALLEL_JOBS=
WAITING_JOB=
WAITING_ON=
CACHE_KEY_FILE=
MASKED_FAILURE_JOB=
SWALLOWED_STEP=
```

- `TRIGGER_EVENTS`: a comma-separated list of the event names under `on:`
  (e.g. `push,workflow_dispatch`).
- `PARALLEL_JOBS`: a comma-separated list of the two job names that have
  no `needs:` at all.
- `WAITING_JOB`: the one job whose `needs:` names exactly one other job.
- `WAITING_ON`: the job named in that `needs:`.
- `CACHE_KEY_FILE`: the exact filename passed to `hashFiles(...)`.
- `MASKED_FAILURE_JOB`: the job with `continue-on-error: true`.
- `SWALLOWED_STEP`: the exact `name:` of the step whose command ends in
  `|| true`.

Job and step names must match `pipeline.yml` exactly as written.
`SWALLOWED_STEP`'s value contains a space — quote it
(`KEY="value with spaces"`) so sourcing the file doesn't word-split it.

Get a shell: `docker exec -it breaklab-cicd-001 bash`
