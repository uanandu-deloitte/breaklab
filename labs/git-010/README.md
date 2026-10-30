# Undo Safely

`/workspace` has three independent repositories, each with something
that needs undoing -- but not the same way.

## repo-staged

`config.yaml` has a change staged with `git add` that should never have
been staged. Get the working tree back to exactly what the last commit
has -- nothing staged, nothing modified.

## repo-private

The last commit ("DEBUG: temporarily disable auth check") was made
locally and has no remote -- nobody else has ever seen it. Remove it
from history entirely, leaving `HEAD` at the commit before it.

## repo-shared

The last commit ("DEBUG: temporarily disable auth check") is the exact
same mistake -- but this repo is a clone of
`/srv/repo-shared-origin.git`, a bare repo standing in for a shared
remote that already has this commit. Undo its *effect* without
rewriting history: add a new commit that reverses it.

## Record what you used

```
STAGED_CHANGE_COMMAND=
PRIVATE_COMMIT_COMMAND=
SHARED_COMMIT_COMMAND=
```

in `/workspace/answers.env`, one `KEY=value` per line -- just the git
subcommand name you used for each situation (e.g. `restore`, `reset`,
or `revert`), no flags or spaces.

Get a shell: `docker exec -it breaklab-git-010 bash`
