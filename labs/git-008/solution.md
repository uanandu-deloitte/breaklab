# Solution: Fetch, Pull, Push

```bash
cd /workspace

git fetch origin
git log main..origin/main --oneline   # shows the teammate's "Add logging helper" commit
git status                             # "Your branch is behind 'origin/main' by 1 commit"

git pull origin main                   # fast-forwards local main to match origin/main

cat <<EOF >> CHANGELOG.md
## Changelog
- Added greeting module
- Pulled in logging helper from teammate
EOF
git add CHANGELOG.md
git commit -m "Add changelog"

git push origin main
```

`git fetch` downloads the teammate's commit into `origin/main` without
touching local `main`, so `git log main..origin/main` and `git status`
can show exactly what's new before anything changes locally. Since local
`main` hasn't diverged, `git pull` fast-forwards it straight up to match
`origin/main` — no merge commit needed. The new local commit on top of
that is then pushed, advancing `origin/main` to include both the
teammate's commit and this one.
