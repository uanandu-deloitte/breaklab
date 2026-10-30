# Live config vs. historical record

The same string can legitimately appear in more than one file for very
different reasons — one occurrence is a bug that needs fixing, another
occurrence is a record of history that is correct exactly as it stands.
Telling them apart matters more than finding every match.

## Finding every occurrence

`grep -R "pattern" /some/dir` recursively searches every file under a
directory and prints every matching line with its file path. This is
the easy part — it tells you *where* a string appears, not which
appearance matters.

## What makes a file "live"

A live configuration file is the one a running process actually opens
and reads right now. An archived release snapshot, a changelog entry, or
a migration note that mentions the same value is not live — it's a
record, deliberately preserved, and a correct one at that. Figuring out
which file is which usually comes down to:

- Does the path look like an active deployment location (`config/`,
  not `releases/<old-version>/...`)?
- Does the file's own purpose explain why the old value should still be
  there (a changelog entry describing a past migration, by definition,
  needs to still say what the *old* value was)?

## Why blind find-and-replace is dangerous

A recursive `sed -i 's/old/new/g'` across an entire directory tree finds
every occurrence and changes all of them, with no judgment about which
ones are live and which are historical. Run across `/opt/app` here, it
would "fix" the live config correctly — and simultaneously corrupt the
changelog, the migration notes, and the archived release snapshot,
silently rewriting history to claim the old release used a hostname it
never actually used. The blast radius of a blind replace is every match
in the tree, not just the one that was actually broken.
