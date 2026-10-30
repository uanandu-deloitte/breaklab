# fetch, pull, and push are three different operations

It's easy to treat "sync with the remote" as one action, but git splits
it into distinct steps that do different things — knowing which one
you're running matters once more than one person is pushing to the same
remote.

## git fetch: look, don't touch

`git fetch` downloads any commits the remote has that you don't, and
updates your remote-tracking ref (`origin/main`) to point at them. It
does **not** touch your own branch (`main`) at all. After a fetch,
`origin/main` and `main` can point at completely different commits —
that's expected, and exactly how you inspect what changed upstream
before deciding what to do about it.

```
git fetch origin
git log main..origin/main --oneline   # commits origin has that you don't
git status                            # "Your branch is behind ... by N commits"
```

## git pull: fetch, then merge

`git pull` is `git fetch` followed by merging `origin/main` into your
current branch (or rebasing, depending on configuration). If your branch
hasn't diverged — you haven't committed anything local since you last
synced — that merge is a fast-forward: your branch's pointer just moves
up to match `origin/main`. If you *have* committed locally, `pull` will
produce a true merge commit (or a conflict, if the same lines changed on
both sides) exactly like merging any other branch would.

## git push: send your commits

`git push origin main` uploads commits your local `main` has that
`origin/main` doesn't, and (assuming the remote hasn't moved further in
the meantime) advances `origin`'s `main` ref to match yours. If someone
else pushed in between your last fetch and your push, git rejects the
push rather than silently overwriting their work — you fetch and merge
(or pull) again first.

## The sequence that keeps everyone in sync

fetch to see what's new -> pull (or merge) to bring it into your branch
-> commit your own work on top -> push to share it. Skipping the
fetch/pull step before pushing is how two people's work silently
conflicts instead of merging cleanly.
