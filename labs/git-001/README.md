# Identify Repo State

You've been handed shell access to `/workspace`, a Git repository
someone else was working in. Before you change anything, read its
current state the same way you'd read an unfamiliar host's identity
first.

Record these values in `/workspace/answers.env`, one `KEY=value` per
line:

```
CURRENT_BRANCH=
HEAD_COMMIT_SHORT_SHA=
OTHER_BRANCH_NAME=
MODIFIED_FILE_COUNT=
```

- `CURRENT_BRANCH` — the branch currently checked out.
- `HEAD_COMMIT_SHORT_SHA` — the short SHA of the commit `HEAD` points at.
- `OTHER_BRANCH_NAME` — the name of a branch that exists in this repo
  but is not the one checked out.
- `MODIFIED_FILE_COUNT` — how many tracked files have uncommitted
  modifications (don't count untracked files).

Collect every value from the running repository itself — don't guess.

Get a shell: `docker exec -it breaklab-git-001 bash`
