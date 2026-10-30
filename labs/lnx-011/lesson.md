# Normalizing messy text with a pipeline

Hand-maintained inventories drift: different people type hostnames in
different case, separate them with commas sometimes and newlines other
times, and leave stray whitespace behind. Fixing this by hand doesn't
scale and introduces new mistakes. A short pipeline fixes the whole file
in one deterministic pass.

## The building blocks

- `tr 'A-Z' 'a-z'` transliterates every uppercase character to its
  lowercase equivalent — a case-insensitive hostname becomes directly
  comparable.
- `tr ',' '\n'` turns every comma into a newline, so a mix of
  comma-separated and newline-separated entries becomes one entry per
  line throughout.
- `sed 's/^[ \t]*//; s/[ \t]*$//'` strips leading and trailing spaces or
  tabs from each line — two substitutions, anchored to the start (`^`)
  and end (`$`) of the line.
- `sort -u` sorts lines alphabetically and removes duplicates in a single
  pass — equivalent to `sort | uniq`, but one command instead of two.

## Order matters

Splitting on commas has to happen before delimiter-sensitive steps run,
or a comma-joined entry gets treated as one long field. Lowercasing
before deduplicating is what makes `WEB-01` and `web-01` collapse into
the same line — `sort -u` only merges lines that are byte-for-byte
identical, so normalization has to happen first.

## Why a pipeline, not manual cleanup

A pipeline re-derives the clean list from the raw input every time it
runs — rerun it against tomorrow's messier inventory and it still
produces a correct, deduplicated result. Manually retyping a cleaned-up
list is a one-time snapshot that's already stale the next time the raw
data changes.
