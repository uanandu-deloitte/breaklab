#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${OFFENDING_PID:-}" ] || fail "OFFENDING_PID is empty"
[[ "$OFFENDING_PID" =~ ^[0-9]+$ ]] || fail "OFFENDING_PID should be a numeric PID, got \"$OFFENDING_PID\""

[ -n "${OFFENDING_PPID:-}" ] || fail "OFFENDING_PPID is empty"
[[ "$OFFENDING_PPID" =~ ^[0-9]+$ ]] || fail "OFFENDING_PPID should be numeric, got \"$OFFENDING_PPID\""
[ "$OFFENDING_PPID" = "1" ] || fail "OFFENDING_PPID should be 1 — all three background processes were forked directly from the container's init"

[ -n "${OFFENDING_USER:-}" ] || fail "OFFENDING_USER is empty"
[ "$OFFENDING_USER" = "root" ] || fail "OFFENDING_USER should be root — that's who started all three background processes"

[ -n "${OFFENDING_CMD:-}" ] || fail "OFFENDING_CMD is empty"
echo "$OFFENDING_CMD" | grep -qiE 'compactor|while true' || \
  fail "OFFENDING_CMD (\"$OFFENDING_CMD\") doesn't look like the actual offending process — identify the one burning CPU, not a healthy worker"

# The runaway busy loop must no longer be running.
if pgrep -f 'while true; do :; done' >/dev/null; then
  fail "the runaway process is still running — terminate it"
fi

# Both healthy workers must still be running (undisturbed).
pgrep -f 'worker-heartbeat' >/dev/null || fail "worker-heartbeat is no longer running — only the offending process should have been stopped"
pgrep -f 'worker-metrics-sync' >/dev/null || fail "worker-metrics-sync is no longer running — only the offending process should have been stopped"

echo "PASS: the runaway process is stopped, both healthy workers are still running, and answers.env is filled in"
