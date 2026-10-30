# Cache Keys Must Track What Changed

This container has a dependency lockfile at `/srv/app/requirements.lock`
and a script, `/usr/local/bin/cache-key`, whose job is to print a key
representing "the current state of the dependencies." A build system
would use that key to decide whether its cache is still good — but
first, figure out whether this key actually does its job.

A frozen, untouched copy of the original script is kept at
`/opt/reference/cache-key.orig` (read-only) so you always have
something to diagnose even after you've started editing the real one.

## Tasks

1. Read `/usr/local/bin/cache-key`. Figure out exactly what it computes
   its key from.
2. Run it (`/usr/local/bin/cache-key`), note the key it prints, then
   edit `/srv/app/requirements.lock` (append a line, change a version —
   anything) and run it again. Does the key change?
3. Fix `/usr/local/bin/cache-key` so its key is actually derived from
   `requirements.lock`'s content — a hash (`sha256sum`) is the standard
   approach.
4. Confirm your fix: editing the lockfile should now change the key;
   leaving it alone should produce the same key every time.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
ORIGINAL_KEY_SOURCE=
CACHE_KEY_NOW_TRACKS_LOCKFILE=
```

- `ORIGINAL_KEY_SOURCE`: a short phrase naming what the *original*
  (unfixed) `cache-key` actually computed its key from — check
  `/opt/reference/cache-key.orig` for this, since you'll have already
  edited the real one.
- `CACHE_KEY_NOW_TRACKS_LOCKFILE`: `yes` or `no` — after your fix, does
  editing `requirements.lock` actually change what `cache-key` prints?

Get a shell: `docker exec -it breaklab-cicd-006 bash`
