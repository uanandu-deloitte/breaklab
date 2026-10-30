# A container is not a tiny virtual machine

It's easy to picture a container as a little computer with its own kernel,
its own process table, and its own name. None of that is true. A container
is a regular process (or group of processes) on the host, made to *look*
isolated by a handful of Linux kernel features — namespaces for what it can
see, cgroups for what it can use. Nothing about that process gets its own
kernel, and nothing about its isolation is physical.

## What's actually isolated, and what isn't

- **Process view**: a PID namespace gives the container its own numbering
  starting at PID 1, and `ps` inside the container can only see processes in
  that namespace — not the host's, and not any other container's. This is
  why `ps aux` inside a freshly started container shows a handful of
  entries, not the hundreds a real host would have.
- **The kernel**: there is exactly one kernel involved — the host's. A
  container image bundles a userspace (libraries, binaries, an `/etc`) but
  never a kernel. `uname -r` inside any container on a given host always
  reports that host's kernel version, regardless of which image it's
  running — `ubuntu:24.04`, `alpine:3.19`, anything.
- **Hostname**: the Dockerfile never sets this. Unless a human overrides it
  at run time (`docker run --hostname ...`), the container runtime assigns
  the container's own short container ID as its hostname the moment it
  starts — it's runtime-assigned identity, not build-time identity.

## Why this matters

Confusing "isolated" with "independent" leads to real mistakes: assuming a
kernel upgrade inside a container is possible (it isn't — you'd have to
reboot the host), assuming two containers on the same host can't see or
affect each other's resource usage (cgroups limit this, but the kernel is
still shared and host-level exhaustion affects everyone), or assuming a
container's hostname means anything about its image. The isolation is real
and useful, but it is a view, not a separate machine.
