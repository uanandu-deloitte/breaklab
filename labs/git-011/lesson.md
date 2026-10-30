# Stash: a side pocket for work in progress

Sometimes you need a clean working tree *right now* -- to switch
branches, pull a fix, or just take stock -- but you're not ready to
commit what you're in the middle of. `git stash` is exactly that: it
lifts your uncommitted changes off the working tree and parks them on a
separate stack, leaving the tree clean without committing anything.

## What gets stashed, by default

Plain `git stash` (or `git stash push`) captures:

- changes to tracked files, staged or not

It does **not** capture untracked files -- new files git doesn't know
about yet -- unless you ask it to.

## Including untracked files

`git stash push -u` (`--include-untracked`) sweeps up new, untracked
files along with the usual tracked changes, so a single stash entry
really does capture everything sitting in the working tree.

## Getting it back

- `git stash pop` re-applies the most recent stash entry *and* removes
  it from the stack in one step -- the normal move once you're done
  with it.
- `git stash apply` re-applies it but leaves it on the stack -- useful
  if you want to apply the same stash in more than one place, but easy
  to forget about afterward if that's not actually what you meant to
  do.

## Checking your work

`git stash list` shows what's still parked. An empty list after a
`pop` confirms the stash was actually consumed, not just applied and
left behind as a duplicate of work that's now also back in your
working tree.
