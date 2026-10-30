# A pipeline file, read like a map

A CI workflow file describes four things at once: what makes it run, what
order its jobs run in, what it reuses between runs, and what it does when
something fails. All four are readable directly from the YAML, before any
pipeline ever executes.

## Where to look

- **Triggers** live under `on:`. A workflow can have more than one — a
  `push` trigger and a manual `workflow_dispatch` trigger both firing the
  same workflow is normal, not a mistake.
- **Job ordering** is controlled by `needs:`. A job with no `needs:` key
  starts as soon as the workflow does — if two jobs both have no `needs:`,
  they run in parallel. A job only waits for another when it names that
  job under its own `needs:`.
- **Caching** is usually keyed with `hashFiles('some-lockfile')`. That
  hash becomes part of the cache `key:`; change the lockfile's contents
  and the hash changes, which changes the key, which means the next run
  can no longer find a matching cache entry.
- **Swallowed failures** hide in two places: a job marked
  `continue-on-error: true` reports its own failure without stopping
  anything that depends on it, and a shell step ending in `|| true` turns
  that step's real exit code into `0` no matter what actually happened.

## Why this matters

None of this requires running the pipeline to find out. A workflow file
is a complete, static description of behavior — the same discipline that
reads a Dockerfile or a shell script to predict what it does, applied to
CI configuration instead.
