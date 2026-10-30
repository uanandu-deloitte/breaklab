# Tag a Release, Cherry-Pick a Fix

`/workspace` has a `release` branch, tagged `v1.0.0` (an annotated
tag), and a `main` branch that's moved ahead with three more commits
since release branched off.

Exactly one of those three commits is a hotfix that needs to go out on
`release` too: **the commit that fixes the null-pointer crash in
`handler.py`**. The other two commits on `main` are unrelated and must
not come along.

## Task

1. Find that commit on `main` by its content/message -- not by
   guessing or copying a SHA from here, there isn't one.
2. Cherry-pick only that commit onto `release`.
3. Create a new annotated tag, `v1.0.1`, at the result.

## Constraints

- `release` must end up with the null-pointer fix and *nothing else*
  from `main`'s other two commits.
- `v1.0.1` must be a real annotated tag (`git tag -a`), not a
  lightweight one.

Get a shell: `docker exec -it breaklab-git-012 bash`
