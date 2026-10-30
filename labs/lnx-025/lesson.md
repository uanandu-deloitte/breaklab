# Why "it ran without errors" isn't the same as "it worked"

A shell script with no error-handling can run to completion, print
"success," and exit `0` while having done something completely different
from what was intended. None of that is a syntax error — it's just a
script that was never taught to check its own work.

## Quoting is not optional

An unquoted variable expansion is subject to word-splitting on
whitespace before the command ever sees it. `cp $f /dest/` with
`f="report draft.txt"` doesn't pass one argument — it passes `report`
and `draft.txt` as two separate (and individually nonexistent) paths.
Quoting (`cp "$f" /dest/`) is what guarantees a variable's value reaches
the command as a single argument, no matter what it contains.

## Checking exit codes, not assuming them

Every external command (`cp`, `tar`, `mkdir`, ...) returns an exit code
that says whether it actually succeeded. `command || { ...; exit 1; }`
(or an explicit `if ! command; then ...; fi`) is how a script finds out.
Without that check, a script's later "success" message is just a guess —
and in a script that ran unattended, a wrong guess is the only
notification anyone gets.

## Only report success on the path that earned it

The fix for "prints success unconditionally" isn't deleting the success
message — it's making sure the only way execution reaches it is by
passing every check along the way. Every failure path should `exit`
before ever reaching that line.
