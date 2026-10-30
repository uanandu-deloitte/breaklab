# Log Rotation Breaks Naive File Tailing

A reader process has been following `/var/log/app.log` since this
container started. At some point last night (actually: a second or two
after this container started), that log got rotated — the old file was
renamed aside and a fresh one took its place, exactly the way a real
`logrotate` job would do it. The reader never restarted and never
complained. Figure out exactly what it's reading now.

## Tasks

1. Find the reader's PID (`/var/run/reader.pid`) and compare the line
   counts of `/var/log/app.log` and `/var/log/reader-output.log` — the
   reader's own output file.
2. List the reader process's open file descriptors under
   `/proc/<pid>/fd/` and resolve each one with `readlink -f` to see
   exactly which file it points at right now.
3. Compare that file's inode against the current `/var/log/app.log`'s
   inode with `stat -c %i` to confirm they are genuinely two different
   files, not just two different names for the same one.

Record your findings in `/workspace/answers.env`, one `KEY=value` per line:

```
READER_STALE=
READER_FD_PATH=
INODES_MATCH=
```

- `READER_STALE`: `yes` or `no` — is the reader's output missing lines
  that are already present in the current `app.log`?
- `READER_FD_PATH`: the exact path (as `readlink -f` resolves it) that
  the reader's log-reading file descriptor currently points at.
- `INODES_MATCH`: `yes` or `no` — does that file have the same inode
  number as the current `/var/log/app.log`?

Get a shell: `docker exec -it breaklab-obs-005 bash`
