# Installing exactly what you need, nothing more

A minimal base image ships with almost nothing beyond a shell and core
utilities — intentionally, so the image stays small and its attack
surface stays small. Every troubleshooting tool you'll actually want
during an incident gets added deliberately, not by accident.

## The two-step install

```bash
apt-get update
apt-get install -y <package> [<package> ...]
```

`apt-get update` refreshes the local package index from the
distribution's repositories — without it, `apt-get install` is working
from a (possibly empty, possibly stale) index and will often fail to
find a package that genuinely exists. This is the single most common
reason a fresh install fails on a brand-new container.

## A package name isn't always the binary name

The package you install and the command you run aren't always spelled
the same way — `iproute2` is the package; `ip` and `ss` are the commands
it provides. When you're not sure which package provides a binary you
need, `apt-cache search <keyword>` or a quick search of the
distribution's package index will usually tell you.

## Verify, don't assume

"Installed" and "working" aren't the same claim. After installing a
tool, run it — its `--version`, `-V`, or `-v` flag is almost always safe
to run immediately and confirms three things at once: the binary exists,
it's on `PATH`, and it actually executes without error.

## Install only what the task calls for

Every package added to a host is something that has to be patched,
tracked, and justified later. Resist the habit of throwing in a handful
of "useful" extras beyond what's actually needed — a minimal, deliberate
set of tools is easier to reason about than a host that's slowly
accumulated everything anyone ever found convenient.
