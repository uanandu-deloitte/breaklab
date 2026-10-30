# Reading a log without opening the whole thing

A log file can be thousands of lines long. You almost never need to read
every one of them — a handful of small, purpose-built commands can answer
most questions about a log faster than scrolling through it.

## The building blocks

- `head -n <N>` prints the first `N` lines (default 10); `tail -n <N>`
  prints the last `N`. `head -1` / `tail -1` get you just the very first
  or last line — useful for "when did this log start/end?"
- `wc -l` counts lines. Piped after a filter, it counts how many lines
  matched.
- `cut -d' ' -f<N>` pulls out column `N` of a space-delimited line — given
  a log's fixed column layout, this isolates just the field you care
  about (a status code, a path, an IP).
- `sort` groups identical lines together; `uniq -c` then counts how many
  consecutive identical lines there are. `sort | uniq -c | sort -rn` is
  the standard idiom for "rank values by how often they appear."

## Combining them

To count how many requests hit a given status class, filter the status
column first (`grep` or `awk` on the right field), then pipe into `wc -l`.
To find the most frequent value in a column — a path, an IP, anything —
extract that column, then run it through `sort | uniq -c | sort -rn` and
look at the top line.

## Why extraction, not inspection

Reading the log with your eyes to answer "which path is most common?"
doesn't scale past a dozen lines and is easy to get wrong by missing a
line. Piping the real column through `sort | uniq -c` gives you an exact
count, and the same pipeline still works correctly if the log were ten
thousand lines instead of a dozen.
