# Solution: Tag a Release, Cherry-Pick a Fix

```bash
cd /workspace
git log --oneline main           # find the hotfix commit by its message
# ... abc1234 Update changelog for v1.1
# ... def5678 Fix null-pointer crash in handler.py   <-- this one
# ... 9abcdef Add rate limiting to API

git checkout release
git cherry-pick def5678          # use whatever SHA your own log showed
git tag -a v1.0.1 -m "Release 1.0.1: null-pointer crash fix"
```

`git log --oneline main` is enough to identify the right commit by its
message and the problem it describes -- no need to inspect every diff.
`git cherry-pick` replays only that one commit's change onto `release`,
leaving the other two main-only commits (and everything that came
after them) behind entirely. The new annotated tag records the hotfix
release without ever merging `main` into `release` wholesale, which
would have dragged the rate-limiting and changelog commits along with
it.

This exercise cherry-picks straight onto `release` rather than a new
branch off it, since the task is "ship a point release," not "prepare
something for further review" -- tagging the result directly is the
simpler, equally correct path.
