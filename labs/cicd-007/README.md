# Promoting Artifacts Between Environments

`/srv/artifacts/` has three build artifacts in it, plus a `current`
entry. Your job is to figure out which artifact is actually marked for
promotion right now — without assuming the newest-looking filename is
the answer.

## Tasks

1. List the artifacts directory: `ls -la /srv/artifacts/`. Notice that
   `current` isn't a regular file.
2. Find out exactly what `current` points to: `readlink
   /srv/artifacts/current`.
3. Compare that against which file has the highest version number in
   its name, and which file was built most recently
   (`ls -la --time-style=full-iso /srv/artifacts/*.tar`). Are they the
   same file as the symlink's target?
4. Decide: if you'd guessed "promoted" based on filename or recency
   alone, would you have gotten it right?

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
PROMOTED_ARTIFACT=
WOULD_GUESSING_BY_NAME_BE_RIGHT=
```

- `PROMOTED_ARTIFACT`: the exact filename `/srv/artifacts/current`
  actually points to right now.
- `WOULD_GUESSING_BY_NAME_BE_RIGHT`: `yes` or `no` — would picking the
  highest-version or most-recently-built filename have given you the
  same answer as the symlink?

Get a shell: `docker exec -it breaklab-cicd-007 bash`
