# The Staging Area

`/workspace` is a git repo with two things sitting in its working tree
right now:

- `config.yaml` has an uncommitted edit: its `timeout` was changed from
  `10` to `30`.
- `scratch-notes.txt` is a brand new file that was never committed —
  local notes that shouldn't go into the repo.

## Tasks

1. Use `git status`, `git diff`, and `git diff --staged` to confirm
   which file is modified-but-tracked and which is untracked.
2. Stage **only** `config.yaml`'s change.
3. Commit it with exactly this message: `Increase API timeout to 30s`
4. Leave `scratch-notes.txt` untracked — don't stage or commit it.

Get a shell: `docker exec -it breaklab-git-002 bash`
