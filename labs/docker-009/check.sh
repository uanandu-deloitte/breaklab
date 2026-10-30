#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

# Determine the real cgroup version and the real paths for each ceiling
# dynamically, rather than assuming one -- the whole point of this lab.
if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
  real_version=v2
  pids_path=/sys/fs/cgroup/pids.max
  mem_path=/sys/fs/cgroup/memory.max
  cpu_path=/sys/fs/cgroup/cpu.max
else
  real_version=v1
  pids_path=/sys/fs/cgroup/pids/pids.max
  mem_path=/sys/fs/cgroup/memory/memory.limit_in_bytes
  cpu_path=/sys/fs/cgroup/cpu/cpu.cfs_quota_us
fi

[ -n "${CGROUP_VERSION:-}" ] || fail "CGROUP_VERSION is empty"
got_version=$(echo "$CGROUP_VERSION" | tr '[:upper:]' '[:lower:]')
[ "$got_version" = "$real_version" ] || fail "CGROUP_VERSION ($CGROUP_VERSION) doesn't match what's actually mounted here ($real_version)"

[ -f "$pids_path" ] || fail "expected cgroup pids file ($pids_path) not found for cgroup $real_version"
real_pids_max=$(cat "$pids_path")
[ -n "${PIDS_MAX:-}" ] || fail "PIDS_MAX is empty"
[ "$PIDS_MAX" = "$real_pids_max" ] || fail "PIDS_MAX ($PIDS_MAX) doesn't match the live cgroup file ($real_pids_max)"

real_mem_max=$(cat "$mem_path" 2>/dev/null || echo "unreadable")
[ -n "${MEMORY_MAX:-}" ] || fail "MEMORY_MAX is empty"
[ "$MEMORY_MAX" = "$real_mem_max" ] || fail "MEMORY_MAX ($MEMORY_MAX) doesn't match the live cgroup file ($real_mem_max)"

real_cpu_max=$(cat "$cpu_path" 2>/dev/null || echo "unreadable")
[ -n "${CPU_MAX:-}" ] || fail "CPU_MAX is empty"
[ "$CPU_MAX" = "$real_cpu_max" ] || fail "CPU_MAX ($CPU_MAX) doesn't match the live cgroup file ($real_cpu_max)"

# Ground truth: this lab sets a real pids_limit, so the live ceiling must
# be an actual number, not "max"/unlimited -- otherwise the fixture itself
# is broken, independent of what the learner answered.
case "$real_pids_max" in
  max|unreadable) fail "this container's pids ceiling reads as unlimited ($real_pids_max) -- expected a real number from this lab's pids_limit" ;;
esac

echo "PASS: answers.env correctly reports this container's real cgroup version ($real_version) and its pids/memory/CPU ceilings"
