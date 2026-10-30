# Permission Playground

`/workspace/playground/` has three files, each with different ownership
and permissions, and this host also has an unprivileged user named
`tester`.

```
owner-only.txt      owned root:root,    mode 600
group-readable.txt  owned root:opsgroup, mode 640
world-readable.txt  owned root:root,    mode 644
```

`tester` is not a member of `opsgroup`.

## Tasks

For each file: predict whether `tester` can read it, then verify your
prediction by actually testing access as `tester` (not as root).

Write your verified answers to `/workspace/answers.env`:

```
OWNER_ONLY_READABLE_BY_TESTER=
GROUP_READABLE_READABLE_BY_TESTER=
WORLD_READABLE_READABLE_BY_TESTER=
```

Each value must be exactly `yes` or `no`.

## Constraints

- Don't change any file's permissions or ownership — this lab is about
  reading the existing state correctly, not fixing anything.
- Base every answer on an actual test run as `tester`, not a guess.

Get a shell: `docker exec -it breaklab-lnx-014 bash`
