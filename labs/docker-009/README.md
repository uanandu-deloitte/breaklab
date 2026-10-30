# Cgroup Ceilings: PIDs, Memory, CPU

This container was started with a real `--pids-limit`. Nothing about that
limit shows up in `ps` or `top` — find the real ceiling, and the real
ceilings for memory and CPU alongside it, directly from the cgroup
filesystem.

## Tasks

1. Run `mount | grep cgroup` and decide whether this container sees
   cgroup v1 (several separate controller mounts) or cgroup v2 (one
   unified mount).
2. Run `ps aux | wc -l` and confirm it tells you nothing about any
   ceiling — only a count.
3. Using the right path for whichever cgroup version this host actually
   mounts, read the real ceiling and current usage for:
   - PIDs (`pids.max` / `pids.limit`, and `pids.current` / `pids.usage`)
   - Memory (`memory.max` / `memory.limit_in_bytes`)
   - CPU (`cpu.max` / `cpu.cfs_quota_us`)

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
CGROUP_VERSION=
PIDS_MAX=
MEMORY_MAX=
CPU_MAX=
```

- `CGROUP_VERSION`: `v1` or `v2` — whichever this container actually has
  mounted.
- `PIDS_MAX`: the exact value read from this container's real pids
  ceiling file.
- `MEMORY_MAX`: the exact value read from this container's real memory
  ceiling file.
- `CPU_MAX`: the exact value read from this container's real CPU ceiling
  file.

Collect every value from the running container itself — don't guess
which cgroup version is in play, check first.

`CPU_MAX`'s value contains a space on cgroup v2 (`cpu.max` reads back as
two fields, e.g. `max 100000`) — quote it (`CPU_MAX="value with spaces"`)
so sourcing the file doesn't word-split it.

Get a shell: `docker exec -it breaklab-docker-009 bash`
