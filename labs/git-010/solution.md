# Solution: Undo Safely

```bash
cd /workspace/repo-staged
git restore --staged config.yaml
git restore config.yaml

cd /workspace/repo-private
git reset --hard HEAD~1

cd /workspace/repo-shared
git revert --no-edit HEAD

cat <<EOF > /workspace/answers.env
STAGED_CHANGE_COMMAND=restore
PRIVATE_COMMIT_COMMAND=reset
SHARED_COMMIT_COMMAND=revert
EOF
```

- **repo-staged**: `git restore --staged` unstages the change; `git
  restore` (no `--staged`) then discards it from the working tree too,
  leaving both the index and the working tree matching the last commit.
- **repo-private**: the bad commit has never been seen outside this
  clone, so `git reset --hard` rewriting local history costs nothing --
  no one else's history needs to agree with yours.
- **repo-shared**: the exact same bad commit, but this one is already
  on a remote other clones may have fetched. `git reset` here would
  make your branch and the shared remote's history disagree the moment
  you tried to push -- anyone who already pulled the bad commit would
  be out of sync, and a force-push would be needed to "fix" it
  (destroying their copy of history along the way). `git revert`
  instead adds a new commit on top that undoes the change -- the
  mistake stays visible in history (which is honest), but its effect is
  gone.
