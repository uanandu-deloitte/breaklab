# Cherry-pick: taking one commit without the rest of its branch

Merging or rebasing brings an entire range of commits across at once.
Sometimes you want exactly one of them -- a hotfix that landed on
`main` alongside unrelated work, but needs to ship on a release branch
by itself, without that unrelated work along for the ride.

## What cherry-pick actually does

`git cherry-pick <commit>` takes the *change* a single commit
introduces -- its diff -- and replays it as a brand new commit on your
current branch. The new commit has a new SHA and a new parent; it
isn't the same commit object living in two places, it's a copy of the
change.

## Finding the right commit without a SHA in hand

In real incidents nobody hands you a commit hash. `git log --oneline`
(or `git log -p` if the message alone isn't enough) over the branch
that has the commit you need is usually enough -- look for the
message, or the file, that matches the problem you're actually trying
to fix.

## Why not just merge main into release

Merging (or rebasing release onto main) brings *every* commit main has
accumulated since the branches diverged -- including work that was
never meant to ship in this release. Cherry-pick is the tool for "I
want this one change and nothing else that happened alongside it."

## Tagging the result

An annotated tag (`git tag -a <name> -m <message>`) is a real git
object with its own message, tagger, and date -- unlike a lightweight
tag, which is just a name pointing at a commit. Release tags are almost
always annotated: the message is exactly where "what shipped and why"
gets recorded.
