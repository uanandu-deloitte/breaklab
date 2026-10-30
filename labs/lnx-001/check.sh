#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${DISTRIBUTION:-}" ] || fail "DISTRIBUTION is empty"
os_id=$(. /etc/os-release; echo "$ID")
echo "$DISTRIBUTION" | grep -qi "$os_id" || fail "DISTRIBUTION doesn't mention $os_id"

[ -n "${KERNEL:-}" ] || fail "KERNEL is empty"
[ "$KERNEL" = "$(uname -r)" ] || fail "KERNEL doesn't match uname -r"

[ -n "${HOSTNAME:-}" ] || fail "HOSTNAME is empty"
[ "$HOSTNAME" = "$(hostname)" ] || fail "HOSTNAME doesn't match hostname"

[ -n "${USER:-}" ] || fail "USER is empty"
[ "$USER" = "$(whoami)" ] || fail "USER doesn't match whoami"

[ -n "${SHELL:-}" ] || fail "SHELL is empty"
echo "$SHELL" | grep -qi "bash" || fail "SHELL doesn't look like this container's bash"

[ -n "${WORKING_DIRECTORY:-}" ] || fail "WORKING_DIRECTORY is empty"
[ -d "$WORKING_DIRECTORY" ] || fail "WORKING_DIRECTORY is not a real directory"

[ -n "${ARCHITECTURE:-}" ] || fail "ARCHITECTURE is empty"
[ "$ARCHITECTURE" = "$(uname -m)" ] || fail "ARCHITECTURE doesn't match uname -m"

echo "PASS: answers.env matches the running environment"
