# Finding and clearing disk usage safely

A "disk full" incident is really two separate problems: finding where the
space actually went, and removing only what's genuinely safe to remove.
Conflating the two — deleting the first large thing you find — is how
real data gets lost.

## Narrowing down with du and df

`df -h` shows free space per filesystem, confirming the symptom. `du -sh`
on a directory's immediate children (`du -sh /var/log/*`) shows which
subtree is actually responsible, without having to read through every
individual file. Repeat `du -sh` one level deeper on whichever subtree is
largest until you reach the specific directory that's consuming the
space.

## Disposable vs. retained

Not everything that's large is safe to delete. Generated, regenerable
data — debug output, caches, temp files — is disposable. Data a system
explicitly keeps around on purpose — audit logs, customer records,
application binaries — is retained, regardless of its size. The fix for
a disk-full incident is deleting the disposable data, confirming the
retained data's content is unchanged, and confirming the application can
write again (the `write` didn't just error outright before; it failed
*intermittently*, which is exactly the symptom of a filesystem running
out of free space under load).

## Confirming the fix actually works

A removed junk directory isn't the same thing as a restored service.
Check that the specific capability that failed — "the application
writes new logs" — works again, by actually attempting it (e.g.
`touch`ing a new file), not just by confirming free space increased.
