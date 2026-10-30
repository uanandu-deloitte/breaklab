# Comparing two versions of the same file

Before any config change ships, you need to know exactly what it
changes — not what it was intended to change, what the actual bytes on
disk say will change.

## diff

`diff file1 file2` prints only the lines that differ between two files.
A line starting with `<` is unique to the first file; a line starting
with `>` is unique to the second. Lines that match in both files don't
appear in the output at all — `diff` shows you a delta, not the whole
file twice.

For simple `key: value` config files, each differing line usually means
one key's value changed, which makes `diff`'s output a direct list of
what to review before deploying.

## Why not just eyeball both files

For a handful of lines, reading both files side by side can work. For
anything longer, eyeballing reliably misses small changes — a value
that differs by one digit, or one word, is easy to skim past. `diff`
never misses a differing line, no matter how subtle the change, because
it compares byte by byte rather than by impression.

## Isolating what actually matters

A real config file often has far more lines that stayed the same than
lines that changed. `diff`'s output size is proportional to how much
actually changed, not to the file's total size — which is exactly what
makes it useful for reviewing a change before it goes live, rather than
re-reading an entire file to find what's different.
