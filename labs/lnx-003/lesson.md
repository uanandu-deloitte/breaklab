# The standard directory tree

Every Linux distribution lays out its root directory the same broad way,
which is what lets troubleshooting advice generalize across hosts you've
never logged into before.

## The directories that matter most

- `/etc` holds system-wide configuration — nearly everything you'll edit to
  change how a service behaves lives under here.
- `/var/log` holds logs that change over time: application output,
  authentication attempts, package manager history.
- `/opt` holds optional, self-contained software that doesn't fit the
  distribution's own package layout — a common home for an organization's
  own applications.
- Your home directory (`/root` for root, `/home/<user>` otherwise) holds
  your own files, shell configuration, and nothing the system itself
  depends on.

## Absolute vs. relative paths

An absolute path always starts with `/` and means the same thing no matter
where you currently are. A relative path (`logs/app.log`, `../config`) is
resolved against your current working directory — convenient, but only
correct if you know where you're standing. `pwd` always tells you.
