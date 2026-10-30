# Solution: Stash Your Work

```bash
cd /workspace
git stash push -u -m "wip"
git status            # clean -- nothing staged, modified, or untracked

git stash pop
git status            # app.py modified again, notes.txt back as untracked
git stash list         # empty

cat <<EOF > /workspace/answers.env
STASH_COMMAND=push
STASH_INCLUDE_UNTRACKED_FLAG=-u
POP_COMMAND=pop
EOF
```

`git stash push -u` stashes both the tracked modification to `app.py`
and the untracked `notes.txt` -- plain `git stash push` (no `-u`) would
have left `notes.txt` sitting in the working tree, since untracked
files aren't included by default. `git stash pop` re-applies the single
stash entry and removes it from the stack in the same step -- unlike
`git stash apply`, which would re-apply it but leave a copy sitting on
the stack.
