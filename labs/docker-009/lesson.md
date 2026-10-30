# A resource ceiling doesn't show up where you'd expect

`ps`, `top`, and `docker stats` all show **usage**: how many processes
exist right now, how much memory is actually mapped, how much CPU time
was actually consumed. None of them show a **ceiling** — the hard limit a
container was started with, past which the kernel starts refusing new
allocations outright. A container can be sitting comfortably under every
usage number you'd normally check and still be one allocation away from
hitting a limit none of those tools will ever mention.

This matters most for the process-count limit (`--pids-limit`): once a
container's live process count equals its ceiling, `fork()` itself starts
failing for everything in that container — including the diagnostic
commands you'd reach for to investigate, since running any new command at
all requires a fresh fork/exec. The same shape of problem applies to a
memory limit and a CPU quota: each is a real ceiling enforced by the
kernel's cgroup machinery, and none of the three is visible as a ceiling
in any tool that only reports current usage.

## Where the real ceiling lives

Every resource ceiling a container runs under is written into a cgroup
file, readable directly from inside the container. The catch is that the
exact path depends on which cgroup version this host mounts:

- **cgroup v2** (the modern default on most current Linux hosts) exposes
  one unified hierarchy. A single file, `/sys/fs/cgroup/cgroup.controllers`,
  exists at the root — its presence is the signal you're looking at v2.
  Ceilings live directly under `/sys/fs/cgroup/`: `pids.max`,
  `memory.max`, `cpu.max`.
- **cgroup v1** splits each resource into its own controller directory.
  There is no single `cgroup.controllers` file at the root. Ceilings live
  one directory deeper, under a different filename per resource:
  `/sys/fs/cgroup/pids/pids.max`,
  `/sys/fs/cgroup/memory/memory.limit_in_bytes`,
  `/sys/fs/cgroup/cpu/cpu.cfs_quota_us`.

Checking which version is actually mounted — rather than assuming one
path works everywhere — is the first step, every time. A number read
from the wrong path (or a file that doesn't exist at all) tells you
nothing; the same number read from the right path is the literal ceiling
the kernel is enforcing right now.
