# A secret compared against itself isn't a check

A credential check exists to prove the caller knows something the script
itself doesn't already have memorized. The moment both the submitted value
and the expected value live in the same file, that proof collapses — the
script will always agree with itself, regardless of who's running it or
whether they know anything real. This is the quiet failure mode of
hardcoded secrets in automation: the code looks like it's enforcing access
control, but it's actually just comparing a constant to itself.

## Why this is worse than it looks

A hardcoded secret isn't just bad style — it's a standing leak. Anyone who
can read the script (anyone with access to the repo, the built image, a
backup, a log that happens to print the script's own source) now has the
real credential, forever, until it's rotated everywhere the hardcoded copy
might have spread.

## The fix: read it, don't store it

The script should take the secret as input at the moment it runs —
typically an environment variable injected by whatever orchestrates the
pipeline (a CI system's secret store, a deploy tool, an operator's own
shell) — and compare it against a reference the script can check without
ever holding the plaintext secret itself. A hash of the expected value is
enough to verify a match without the comparison value ever being readable
on disk.

## What "fixed" actually means here

Moving the real secret value from one hardcoded spot to another — even if
it's now compared against an environment variable read — doesn't fix
anything; the real secret is still sitting in a file inside the image. A
genuine fix never bakes the real value in anywhere in the pipeline's own
files, only supplies it at run time.
