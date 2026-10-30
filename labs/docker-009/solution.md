# Solution: Cgroup Ceilings: PIDs, Memory, CPU

```bash
mount | grep cgroup   # tells you v1 (several mounts) vs v2 (one unified mount)

if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
  version=v2
  pids_max=$(cat /sys/fs/cgroup/pids.max)
  mem_max=$(cat /sys/fs/cgroup/memory.max)
  cpu_max=$(cat /sys/fs/cgroup/cpu.max)
else
  version=v1
  pids_max=$(cat /sys/fs/cgroup/pids/pids.max)
  mem_max=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)
  cpu_max=$(cat /sys/fs/cgroup/cpu/cpu.cfs_quota_us)
fi

cat <<EOF > /workspace/answers.env
CGROUP_VERSION=$version
PIDS_MAX=$pids_max
MEMORY_MAX=$mem_max
CPU_MAX="$cpu_max"
EOF
```

`cpu.max` (cgroup v2) reads back as two space-separated fields — quota and
period, e.g. `max 100000` — not a single number. `CPU_MAX` has to be
quoted when it's written, the same way a value with spaces always does in
this lab format, or sourcing `answers.env` later would word-split it and
silently lose everything after the first space.

`ps aux | wc -l` and `top` only ever report how many processes exist and
how much is currently used — neither shows the ceiling this container was
actually started with. That ceiling only exists in a cgroup file, and
which path is real depends on which cgroup version this host mounts:
cgroup v2 keeps one unified hierarchy (`cgroup.controllers` present at the
root), cgroup v1 splits every resource into its own controller directory
with different filenames. `pids.max` reads back this lab's real
`pids_limit` as an actual number; memory and CPU read back `max` (or an
equivalent "unlimited" value for v1) because this lab never set a limit
for those two — the same files are exactly where you'd find a real one if
it had been set.
