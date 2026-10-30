# Shipping logs means tracking what you've already sent

A log shipper's whole job is to move new lines from a source file to a
destination, and nothing else — but "new" only means something if the
shipper remembers what it already sent last time. A naive copy-the-whole
-file approach works once and then duplicates everything on every
subsequent run.

## Why this matters

Most real shippers tail a growing source file forever: an application
keeps writing to it, and the shipper runs again and again, each time
picking up only the lines appended since its last pass. Getting this
wrong either loses lines (if the tracked offset jumps too far ahead) or
duplicates them forever (if it never moves at all).

## Where to look

- `/var/log/app.log` is the source — a file that keeps growing.
- `/var/log/aggregator.log` is the destination.
- `ship-logs --source <path> --destination <path>` runs one shipping
  pass: it figures out how many lines it already shipped from that
  source last time, and appends only the lines after that point.

## Proving it actually works

Running `ship-logs` once should make the destination match the source
exactly. The real test is running it a *second* time after the source
has grown — if the destination still matches the source exactly (not
double the original lines, and not missing the new one), the offset
tracking is genuinely working.
