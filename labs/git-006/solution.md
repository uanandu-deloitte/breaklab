# Solution: Merge It In

```bash
cd /workspace
git checkout main
git merge feature/signup-form
```

`main` has its own "Update README with setup instructions" commit since
the branch point, and `feature/signup-form` has two commits of its own
("Add signup form template", "Add signup form validation") — neither
branch is an ancestor of the other, so `git merge` can't just move a
pointer forward. It combines both histories into a new commit with two
parents: `main`'s previous tip and `feature/signup-form`'s tip.

No `--no-ff` flag is needed here — the merge *is* a true merge commit
purely because of how the history diverged, not because of any flag.
`git log --graph --oneline --all` afterward shows the two lines
reconverging at the new merge commit.
