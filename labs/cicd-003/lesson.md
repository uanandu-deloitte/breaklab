# An artifact is a contract between stages

A build stage and a deploy stage don't share memory or state — the only
thing connecting them is whatever file the build stage leaves behind and
the deploy stage goes looking for. That file's name, format, and location
are an implicit contract between the two stages, and nothing enforces it
automatically.

## Why this breaks quietly

When a build step's output naming changes — a version suffix gets added,
a directory moves, a file extension changes — nothing about that change is
inherently an error. The build still succeeds, with a clean exit and no
warning. The problem only shows up one stage later, when deploy goes
looking for a filename that no longer exists, because nobody updated the
consumer side of the contract.

## Fixing it durably, not just once

Renaming today's artifact to match what deploy.sh expects "fixes" this one
run, but the same mismatch reappears the next time build.sh's naming
changes — and it will, since artifact names routinely carry version
numbers, build numbers, or commit SHAs that differ on every run. The
durable fix is making deploy.sh *discover* whatever build.sh actually
produced (a glob over the expected prefix, a manifest file, a "latest"
pointer) rather than hunting for one specific filename it has memorized
from a prior release.

## Where to look

- `/opt/pipeline/build.sh` is the one place that truly knows what it just
  built — run it and read its own output.
- `/opt/pipeline/deploy.sh` should locate its input from that same source
  of truth, not from a filename baked in from a previous release.
