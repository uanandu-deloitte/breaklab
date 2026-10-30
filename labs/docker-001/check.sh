#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

CEILING=20

# PROCESS_COUNT_SEEN: must be a small number, close to the live count
# (computed dynamically here rather than hardcoded, since base process
# count can vary by a process or two across hosts/runs).
[ -n "${PROCESS_COUNT_SEEN:-}" ] || fail "PROCESS_COUNT_SEEN is empty"
case "$PROCESS_COUNT_SEEN" in
  ''|*[!0-9]*) fail "PROCESS_COUNT_SEEN is not a plain number: $PROCESS_COUNT_SEEN" ;;
esac

actual_count=$(ps aux --no-headers | wc -l | tr -d ' ')
[ "$actual_count" -lt "$CEILING" ] || fail "this container itself shows $actual_count processes — expected well under $CEILING; the check's assumptions may be stale"

diff=$(( PROCESS_COUNT_SEEN - actual_count ))
[ "$diff" -lt 0 ] && diff=$(( -diff ))
[ "$diff" -le 5 ] || fail "PROCESS_COUNT_SEEN ($PROCESS_COUNT_SEEN) is too far from the live count ($actual_count) in this container"
[ "$PROCESS_COUNT_SEEN" -lt "$CEILING" ] || fail "PROCESS_COUNT_SEEN ($PROCESS_COUNT_SEEN) is not small — a real host would show far more than this container ever does"

# KERNEL_SHARED: the container has no kernel of its own.
[ -n "${KERNEL_SHARED:-}" ] || fail "KERNEL_SHARED is empty"
echo "$KERNEL_SHARED" | grep -qi '^yes$' || fail "KERNEL_SHARED should be 'yes' — the kernel this container sees belongs to the host, not the image"

# HOSTNAME_SOURCE: must name the runtime/container id, not the Dockerfile.
[ -n "${HOSTNAME_SOURCE:-}" ] || fail "HOSTNAME_SOURCE is empty"
echo "$HOSTNAME_SOURCE" | grep -qi 'dockerfile' && fail "HOSTNAME_SOURCE names the Dockerfile — nothing in this image's Dockerfile sets a hostname"
echo "$HOSTNAME_SOURCE" | grep -qiE 'runtime|container[-_ ]?id|docker' || fail "HOSTNAME_SOURCE doesn't name the container runtime/container id as the actual source"

# Ground truth: this container's hostname really does look like a
# runtime-assigned short container id (12 lowercase hex characters), and
# /etc/hostname agrees with the hostname command.
live_hostname=$(hostname)
echo "$live_hostname" | grep -qE '^[0-9a-f]{12}$' || fail "hostname ($live_hostname) doesn't look like a runtime-assigned short container id"
[ "$(cat /etc/hostname 2>/dev/null)" = "$live_hostname" ] || fail "/etc/hostname doesn't match hostname"

echo "PASS: answers.env correctly describes what this container shares with the host and what it doesn't"
