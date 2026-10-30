# Three ways to undo, and why they're not interchangeable

Undoing a git mistake always comes down to one question first: **has
anyone else possibly already seen this history?**

## Not committed yet: restore

`git restore` operates on the working tree and the index (the staging
area) -- it never touches commit history, because there's no commit to
touch yet. `git restore --staged <file>` unstages a change; plain
`git restore <file>` discards an unstaged modification. Together they
get you back to exactly what the last commit recorded.

## Committed, but private: reset

`git reset --hard <commit>` moves your current branch pointer backward
and throws away everything after it -- commits and working tree
changes alike. This is completely safe when nobody else has ever
fetched the commits you're discarding: the history simply never existed
as far as anyone else is concerned.

## Committed, and shared: revert

The instant a commit is pushed somewhere others can pull from, the
calculation changes. Someone else's clone may already have it. Deleting
it from your own branch with `reset` doesn't delete it from theirs --
it just makes your two histories disagree, and the only way to force
your branch's version onto the shared remote is a force-push that
discards anyone else's intervening work along the way.

`git revert <commit>` sidesteps the whole problem: instead of erasing
the bad commit, it adds a brand new commit that applies the exact
opposite change. History stays append-only and consistent for everyone
-- the mistake and its correction are both visible, which is normally
exactly what you want on a shared branch anyway: a record of what
happened, not a rewritten story.

## The rule of thumb

> If it's only in your own local branch and nobody could have pulled
> it: `reset` is fine, and often cleaner.
> If there's any chance someone else already has it: `revert`, never
> `reset`.
