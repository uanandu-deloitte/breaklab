# Branches are cheap; main is precious

A Git branch is nothing more than a movable pointer to a commit.
Creating one costs nothing and risks nothing — the whole point is that
you can do work on a new pointer without ever moving `main`'s pointer
until you explicitly decide to.

## Creating and switching

`git checkout -b <name>` (or the newer `git switch -c <name>`) does two
things in one step: it creates a new branch pointing at your current
commit, and it checks it out, so your next commit lands on the new
branch instead of on `main`.

## Why main doesn't move

Once you're on the new branch, every commit you make moves *that*
branch's pointer forward. `main`'s pointer stays exactly where it was —
nothing about committing on another branch touches it. This is true
even if the new branch's commit changes the same file `main` has;
they're simply different versions of that file, each visible only when
its own branch is checked out.

## Switching back

`git checkout main` (or `git switch main`) moves your working tree back
to whatever `main` points at — your new branch's commit doesn't
"follow you back"; it simply isn't checked out anymore. `git log
--oneline` on `main` after switching back should look exactly like it
did before you ever created the other branch, because nothing about
creating or committing on a branch rewrites history anywhere else.
