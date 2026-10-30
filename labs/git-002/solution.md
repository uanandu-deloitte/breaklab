# Solution: The Staging Area

```bash
cd /workspace
git status
git diff
git diff --staged

git add config.yaml
git commit -m "Increase API timeout to 30s"
```

`git status` shows `config.yaml` under unstaged changes and
`scratch-notes.txt` under untracked files — two different categories.
`git add config.yaml` stages only that file's edit; `scratch-notes.txt`
is never named, so it's never touched. The commit lands with exactly
the requested message and contains only the `config.yaml` change —
`scratch-notes.txt` is still sitting in the working tree afterward,
untracked, exactly as it was before.
