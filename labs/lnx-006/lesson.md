# A file's purpose doesn't depend on where it sits

Moving a file doesn't change what it is. A deployment that puts the
right files in the wrong directories is a very common failure mode —
the data survived, but the layout a running application expects to find
it in didn't.

## Reading what a file actually is

- A text config file reads like configuration — `key: value` pairs,
  recognizable settings.
- An executable script usually starts with a shebang line (`#!/...`)
  and its content is clearly a program, not data or settings.
- A data file's content looks like the application's own data format,
  not configuration and not a program.

`cat` (for small files) or `file` tells you what you're actually looking
at, independent of whatever directory it happens to be sitting in right
now or whatever its filename suggests.

## mv vs. permissions

`mv` relocates a file without reading and rewriting its content — fast,
and nothing about the file's bytes changes. What `mv` does *not*
guarantee is that the file's permission bits still make sense in the new
location: a script that needs to execute still needs its `+x` bit
checked with `ls -l` after the move, and restored with `chmod +x` if it
didn't survive however the file got misplaced in the first place.

## Why not just delete and recreate

Recreating a "placeholder" file in the right place is never equivalent
to moving the real one — you'd lose whatever the original file actually
contained. The fix here is purely about location: every byte the
original files held needs to end up intact, just in the right directory.
