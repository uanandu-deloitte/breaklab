# Pipes, redirection, and counting without counting

The shell's real power isn't any single command — it's chaining small,
single-purpose commands together so each one's output becomes the next
one's input.

## The building blocks

- `|` (pipe) sends one command's stdout directly into the next command's
  stdin, with nothing touching disk in between.
- `>` redirects stdout to a file, overwriting it; `>>` appends instead.
- `2>` redirects stderr specifically — useful when you want a command's
  errors out of your way without losing its real output.
- `sort | uniq -c` is the standard idiom for "count how many times each
  distinct value appears": `sort` groups identical lines together so
  `uniq -c` can count consecutive duplicates.
- `sort -rn` after `uniq -c` puts the highest count first, numerically
  (plain `sort` would order "10" before "9" as text).

## Why a pipeline, not a manual count

A pipeline re-derives its answer from the current input every time it
runs. A manually counted answer typed into a file is only ever correct
for the exact data you looked at — reorder the input, add one more line,
and a hard-coded answer silently goes stale while a real pipeline doesn't.
