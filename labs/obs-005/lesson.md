# A file path and an open file descriptor are not the same thing

When a process opens a file, the kernel hands it back a descriptor
pointing at that file's actual data — its inode — not at the path string
used to open it. For as long as that descriptor stays open, it keeps
pointing at the same data, even if the path it was opened through gets
renamed, replaced, or deleted out from under it. This is exactly what
happens during log rotation.

## What rotation actually does

A typical rotation job does roughly this:

```
mv app.log app.log.1
touch app.log
```

The application writing logs reopens or was never holding a long-lived
handle, so it adapts fine (or a rotation-aware tool sends it a signal to
reopen). But a separate reader — a metrics agent, a naive `tail -f`, a
shipper — that already had `app.log` open before the rotation keeps its
original descriptor. That descriptor is still perfectly valid; it just
now belongs to a file reachable only via the name `app.log.1`. Nothing
new ever gets appended to `app.log.1` again, so the reader goes quiet
forever, having no idea the world moved on without it.

## Why this produces no error

Both sides of this are working exactly as designed:

- The application opens `app.log` by path and writes to it successfully.
- The reader reads from its open descriptor successfully.

Neither process does anything wrong. The failure is purely that the path
`app.log` and the reader's descriptor now refer to two different files —
and nothing forces them to notice.

## How to actually see this

- `/proc/<pid>/fd/` lists every file descriptor a process has open, each
  one a symlink you can resolve with `readlink -f` to see exactly which
  file it targets right now, independent of what the path "should" mean.
- `stat -c %i <path>` prints a file's inode number. Two different files
  (even briefly sharing a name, or formerly sharing one) will have
  different inodes; the same file, renamed, keeps the same inode.
- `lsof -p <pid>` is a friendlier front end to the same information.

## The fix

A reader needs to either reopen the file by path periodically (`tail -F`
does exactly this instead of `tail -f`), or the rotation strategy needs
to avoid breaking existing descriptors in the first place (copy the file
then truncate the original in place, rather than renaming it — the
"copytruncate" rotation strategy).
