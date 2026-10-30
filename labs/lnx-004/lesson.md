# Searching by content, not just by name

Finding a file by its name is easy when you already know the name.
Finding *which* file contains a specific piece of data is a different
problem — one that comes up constantly in real incidents, where you know
what you're looking for (an ID, an error string, a hostname) but not
which file recorded it.

## Two tools, two jobs

- `find` locates files by name, location, type, or metadata —
  `find /var/log -name '*.log'` lists every `.log` file under a
  directory, but it never looks inside them.
- `grep` searches inside file content. `grep -r "pattern" /some/dir`
  searches every file under a directory, recursively; adding `-l`
  prints just the matching file's path instead of the matching lines.

## Combining them

`find / -type f -exec grep -l "pattern" {} +` runs `grep -l` against
every regular file `find` locates, batching files into each `grep`
invocation for speed. `grep -R "pattern" /` does nearly the same thing
more directly, recursing through everything from `/` down.

## Why you can't guess the location

Real deployments scatter logs across `/var/log`, an application's own
`/opt/<app>/logs`, and temporary debug dumps left in `/tmp` — and
archived or rotated logs often end up nested several directories deeper
than the active log, under a name that doesn't obviously say "log" at
all. A content search finds the right file regardless of where it ended
up or what it's called; guessing a path does not.
