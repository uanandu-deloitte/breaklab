# umask: the mask applied before a file is born

Every time a process creates a file, it asks the kernel for a specific
starting permission — almost always `666` (`rw-rw-rw-`) for a plain file,
`777` (`rwxrwxrwx`) for a directory. The kernel never grants that request
as-is: it first subtracts whatever bits the creating process's **umask**
says to remove. What's left is the file's actual starting mode.

## Reading a umask

`umask` with no arguments prints the mask currently in effect for your
shell. Every file your shell (or anything it launches) creates from that
point on is filtered through it.

A umask is three octal digits, one per permission class — owner, group,
other — same ordering as `chmod`. Each digit lists the bits to **remove**,
not the bits to keep. A requested `666` filtered through a umask of `027`
removes nothing from owner, removes the write bit from group, and removes
everything from other, leaving `640` (`rw-r-----`).

## Choosing a umask for shared, group-collaborative work

A common professional default for a directory a whole group works in
together is `umask 007`: it removes nothing from owner or group (both
keep full read/write on a new file) and removes everything from other.
Applied to the default file request of `666`, that yields `660`
(`rw-rw----`) — the owner and the group can both read and write the file;
everyone else gets nothing at all.

## Why the default is often wrong for this

A looser umask like the common interactive default of `022` strips only
the write bit from group and other, leaving new files at `644`
(`rw-r--r--`) — readable by literally anyone on the system, which is
exactly the kind of accidental over-exposure a shared team directory
should never produce. The umask is set once, per session, but it silently
shapes the permissions of every file created afterward — getting it wrong
is invisible until someone notices who can read something they shouldn't.
