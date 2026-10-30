# Curate a Commit

`/workspace` is a git repo with two unrelated changes already applied
to the working tree, uncommitted:

- `app.py`: a real bugfix — `get_last_item` was reading one index past
  the end of the list; it's been corrected.
- `utils.py`: an unrelated stray debug `print(...)` statement, left in
  at the top of the file.

Neither change is committed yet, and they have nothing to do with each
other.

## Tasks

- Inspect the working tree with `git status` and `git diff` before
  staging anything.
- Create **two separate commits**: one containing only the `app.py`
  fix, one containing only the `utils.py` change.
- Do not end up with a single commit that touches both files.
- Don't leave anything uncommitted when you're done.

Get a shell: `docker exec -it breaklab-git-003 bash`
