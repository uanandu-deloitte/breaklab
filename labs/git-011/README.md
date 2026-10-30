# Stash Your Work

`/workspace` has two kinds of uncommitted work sitting in the working
tree:

- a modification to the tracked file `app.py`
- a brand new, untracked file `notes.txt`

## Task

1. Stash *everything* -- the tracked change and the untracked file both
   -- in one stash entry.
2. Confirm the working tree is completely clean afterward (`git
   status` should show nothing to commit, nothing untracked).
3. Re-apply the stash, and confirm both the `app.py` change and
   `notes.txt` are back exactly as they were.
4. The stash stack should end up empty -- popped, not left sitting
   there "just in case".

Record what you used in `/workspace/answers.env`, one `KEY=value` per
line (single tokens, no flags combined with spaces):

```
STASH_COMMAND=
STASH_INCLUDE_UNTRACKED_FLAG=
POP_COMMAND=
```

Get a shell: `docker exec -it breaklab-git-011 bash`
