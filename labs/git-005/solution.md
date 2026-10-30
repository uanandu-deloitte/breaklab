# Solution: Wrong Branch

```bash
cd /workspace
git branch feature/csv-export main
git reset --hard known-good-base
```

`git branch feature/csv-export main` creates the new branch pointing at
exactly the commit `main` is on right now — both CSV-export commits are
now reachable from `feature/csv-export` too, not just from `main`.

`git reset --hard known-good-base` then moves `main`'s pointer (and the
working tree) back to the known-good commit. Nothing is lost: the two
commits that used to be "on main" are still full objects in the
repository, still reachable from `feature/csv-export`, which is exactly
what makes this reset safe instead of destructive.
