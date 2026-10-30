# Solution: Curate a Commit

```bash
cd /workspace
git status
git diff

git add app.py
git commit -m "Fix off-by-one in get_last_item"

git add utils.py
git commit -m "Add debug print to utils.py"
```

`git add app.py` stages only the index-off-by-one fix; `utils.py`'s
debug print is a completely different file, so it's left untouched in
the working tree. The first commit's diff is exactly the `app.py` fix,
nothing else. The second `git add utils.py` then stages the remaining
change on its own, so it lands in its own commit instead of ever being
folded into the first one. Had both changes landed in the same file,
`git add -p` would do the same job by letting you pick one hunk at a
time instead of one file at a time.
