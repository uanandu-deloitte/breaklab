# Exit codes are the contract, and pipes can hide them

Every command finishes with an exit code: `0` means success, anything
else means failure. A script that ignores its own commands' exit codes
can print "success" text while its actual work silently failed — the
text and the truth have nothing forcing them to agree.

## Missing preconditions

A script that depends on an environment variable being set needs to check
that explicitly, before using it:

```bash
if [ -z "${APP_CONFIG_PATH:-}" ]; then
  echo "ERROR: APP_CONFIG_PATH is not set; cannot start" >&2
  exit 1
fi
```

Without this check, an unset variable just gets used anyway — often
producing a confusing low-level error from whatever command tried to use
it, instead of a clear message about what the real problem is.

## How a pipe hides a failure

`cmd1 | cmd2` reports the exit status of `cmd2` — the last command in the
pipeline — not `cmd1`. If `cmd1` fails but `cmd2` still runs and succeeds
(or the script never even looks at `$?`), the failure is gone. A script
that always reaches its final `echo "success"` and `exit 0` regardless of
what happened earlier is reporting a status it never actually checked.

## Writing it so failure propagates

- Check required inputs explicitly, with a clear error message, before
  using them.
- Avoid piping a command whose exit status you care about into another
  command, unless you capture and check `${PIPESTATUS[0]}` (or avoid the
  pipe entirely).
- `set -e` makes a script exit immediately on any unchecked command's
  non-zero status — a useful default, but it still won't save a pipeline
  whose last stage happens to succeed while an earlier stage failed.

The only way a startup script can be trusted is if its exit code and its
printed message both reflect what actually happened.
