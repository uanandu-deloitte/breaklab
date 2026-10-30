# "Promoted" is whatever the mechanism says it is — nothing else

Shipping a build to production almost always goes through a
"promotion" step: some mechanism marks one specific artifact as "the
one that's live now." That mechanism might be a config variable naming
an exact file, a symlink pointing at one artifact among several, or a
manifest file recording a version string. Whichever form it takes, it
is the *only* thing that actually defines what's promoted — not the
artifact's filename, not its version number, not when it was built.

## Filenames and timestamps are not the mechanism

It's tempting to look at a directory full of build artifacts and
assume the highest version number, or the most recently built file, is
"obviously" the one that's live. Neither is true in general:

- A team might deliberately pin production to an older, known-good
  build while a newer one bakes in staging — the newer file exists,
  but it was never meant to go live yet.
- A rollback points production back at an older artifact on purpose,
  precisely because the newest one turned out to be bad.

In both cases, the directory still contains a file that looks newer by
every naming convention — and it is still not the one that's
promoted. The only way to know which one actually is comes from
reading the promotion mechanism itself.

## Read the mechanism, not the directory listing

If the mechanism is a symlink, `readlink` tells you its real target —
directly, with no guessing. If it's a config variable, read the
variable. If it's a manifest file, read the manifest. In every case,
the question "what's promoted right now" has exactly one correct
source, and it is never "whichever file sorts last."

## Why this matters

A promotion mechanism that can point at the wrong artifact does so
*silently* — there's no error, because from the tool's point of view
it's just doing exactly what it was told. The only way to catch that
kind of mistake is to already know how to read the real mechanism
yourself, rather than trusting a filename's version number or a
directory listing's sort order to tell you what's actually live.
