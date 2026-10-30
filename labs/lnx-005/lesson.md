# Creating, copying, and moving files

A new application's directory layout usually isn't handed to you fully
built — you create it, then populate it, using a small set of commands
that cover almost every real deployment task.

## The commands

- `mkdir -p` creates a directory and every missing parent directory it
  needs along the way, in one command — without `-p`, `mkdir` refuses to
  create a nested path whose parent doesn't exist yet.
- `touch` creates an empty file, or updates an existing file's timestamp
  without touching its content.
- `cp` copies a file, leaving the original in place. `mv` moves it — the
  original is gone afterward. Confusing the two during a real deployment
  is how source fixtures get accidentally deleted.
- File permissions survive `cp` inconsistently across filesystems; check
  with `ls -l` and re-apply `chmod +x` if a copied script loses its
  executable bit.

## Why "copy, not move" matters here

Treat the original fixture the way you'd treat a release artifact in a
real pipeline: the source of truth stays put, and what you deploy is a
copy of it — so a mistake in the new location never destroys the one
known-good version.
