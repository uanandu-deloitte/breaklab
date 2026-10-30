# One commit, one change

A clean commit history is one where each commit represents a single,
coherent change — a reviewer (including future you) can read one
commit's diff and understand the whole reason for it. A working tree
rarely cooperates with that on its own: real work tends to accumulate
several unrelated edits before any of them get committed.

## Splitting by file

If two unrelated changes land in two different files, splitting them
is direct: stage and commit one file, then stage and commit the other.

```bash
git add file-a
git commit -m "..."
git add file-b
git commit -m "..."
```

`git add <file>` only ever stages that file. Anything you haven't
named stays exactly where it was, so the first commit can't
accidentally pick up the second file's change.

## Splitting within one file

When both unrelated changes land in the *same* file, staging by
filename no longer separates them — `git add file.py` would stage both
at once. `git add -p file.py` solves that by walking through the
file's changes one **hunk** (a contiguous block of changed lines) at a
time, asking `y`/`n`/`s` for each, so you can stage only the hunk that
belongs in the current commit and leave the rest for the next one.

## Why this discipline matters

A commit that mixes a bugfix with an unrelated change is harder to
review, harder to revert cleanly (reverting it undoes both changes,
even though only one was wrong), and harder to `git bisect` against,
since it's no longer obvious which part of the diff caused a
regression. Splitting changes as you commit, rather than after the
fact, is what keeps history actually useful.
