# Why "non-fast-forward" isn't your branch being broken

A `git push` rejection almost always means the same thing: origin has a
commit your local branch doesn't have yet. Git refuses to push *not*
because anything is broken, but because pushing right now would
silently throw away that commit on origin -- overwriting history that
someone else is relying on.

## What actually diverged

Two commits can exist on the same branch name in two different places
at once:

- a commit someone else already pushed to origin, that your local
  clone never fetched
- a commit you made locally, that you haven't pushed yet

Neither commit is wrong. They just both build on the same shared
ancestor, so origin's `main` and your local `main` have quietly become
two different lines of history with the same name.

## Fetch before you decide

`git fetch` updates your knowledge of origin without changing your own
branch -- it's purely informational. Only after fetching can you
compare `origin/main..main` and `main..origin/main` to know exactly
what's on each side before you reconcile them.

## Reconciling without losing anything

- `git pull --rebase` replays your local commit(s) on top of origin's
  updated history -- one straight line, both commits present.
- `git fetch` + `git merge` keeps both histories exactly as they were
  and ties them together with a merge commit -- also both commits
  present, just not a straight line.

What you should *not* do is force-push. `git push --force` makes the
rejection go away by overwriting origin's branch with yours -- which is
exactly the data loss the rejection existed to prevent.
