# Solution: Push Rejected

```bash
cd /workspace
git push origin main
# -> rejected: "Updates were rejected because the remote contains work
#    that you do not have locally."

git fetch origin
git log --oneline main..origin/main   # see the teammate's commit
git log --oneline origin/main..main   # see your own commit

git pull --rebase origin main
git push origin main
```

`git fetch` downloads origin's new commit without touching your own
branch, so you can look at exactly what diverged before deciding how to
reconcile it. `git pull --rebase` replays your local commit on top of
origin's updated `main` instead of creating a merge commit -- the
result is a single straight line of history containing both commits,
and the push that follows is a plain fast-forward.

(A `git fetch` + `git merge origin/main` would also work and is equally
correct here -- it just leaves a merge commit instead of a replayed
one.)
