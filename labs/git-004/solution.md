# Solution: Branch and Switch

```bash
cd /workspace
git checkout -b feature-logging

printf '# Logging plan\n\nAdd structured logging to the service.\n' > logging.md
git add logging.md
git commit -m "Add logging plan"

git checkout main
```

`git checkout -b feature-logging` creates the branch and switches to it
in one step, so the commit that follows lands on `feature-logging`,
not `main`. Switching back with `git checkout main` returns to the
branch pointer that was never touched — `main`'s tip commit and
`git log --oneline` output are identical to what they were before
`feature-logging` ever existed, and `logging.md` is nowhere to be found
there since it only exists on the other branch's commit.
