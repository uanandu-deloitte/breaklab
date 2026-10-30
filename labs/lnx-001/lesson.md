# Linux identity, in six commands

Every Linux host has an identity worth establishing before you change
anything on it: which **distribution** it runs, which **kernel** version,
its **hostname**, which **user** and **shell** you're operating as, and the
system's **architecture**.

## Why this matters

A fix that works on one distribution's package manager, kernel version, or
CPU architecture can silently fail — or silently do the wrong thing — on
another host. Before you run a single troubleshooting command, know what
you're actually troubleshooting.

## Where to look

- `/etc/os-release` holds the distribution name and version in a few
  predictable fields (`ID`, `NAME`, `PRETTY_NAME`).
- `uname -r` prints the running kernel version; `uname -m` prints the CPU
  architecture.
- `hostname` and `whoami` answer "where am I" and "who am I" independently
  of each other.
- `$SHELL` names your configured login shell.
- `pwd` prints your current working directory.

None of these commands need root, and none of them change anything on the
system — this is pure observation before action.
