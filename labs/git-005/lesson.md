# Branches don't have to be created in advance

A feature branch is usually created before you start committing — but
nothing stops you from committing directly on `main` by mistake and
deciding to fix the branching after the fact. Git doesn't care when a
branch is created relative to the commits it ends up pointing at; it
only cares what commit the branch's pointer names right now.

## Moving commits onto a branch they were never on

`git branch <name> <start-point>` creates a new branch pointer at
`<start-point>` without switching to it. If `<start-point>` is `main`'s
current tip, the new branch now points at exactly the same commit `main`
does — the commits that were "mistakenly on main" are now *also*
reachable from the new branch. Nothing moved; a second name was added
for the same history.

## Moving main back without losing anything

Once those commits are reachable from the new branch, `main` no longer
needs to carry them. `git reset --hard <ref>` moves the current branch's
pointer (and the working tree) to `<ref>`, discarding whatever commits
used to be between the old tip and the new one — but "discarding" only
means "no longer on this branch's history." If another branch still
points at or beyond the old tip, those commits are still fully present
in the repository and still reachable.

This is why order matters: branch first, then reset. Reset before
branching and there'd be nothing left pointing at the commits you meant
to save.

## Reading the state before you touch anything

`git log --oneline --all --graph` shows every commit reachable from any
ref, not just the one you're currently on — the fastest way to confirm
where a branch's commits actually are before you move anything.
