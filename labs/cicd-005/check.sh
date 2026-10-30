#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
LOCK=/var/run/deploy/deploy.lock
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
[ -f "$LOCK" ] || fail "$LOCK is missing — fixture was tampered with"

# shellcheck disable=SC1090
source "$ANSWERS"

actual_pid="$(grep -oE '^PID=[0-9]+' "$LOCK" | head -n1 | cut -d= -f2)"
[ -n "$actual_pid" ] || fail "couldn't read a PID out of $LOCK — fixture was tampered with"

# LOCK_PID: must match whatever PID the lock file actually names right
# now, read live rather than assumed.
[ -n "${LOCK_PID:-}" ] || fail "LOCK_PID is empty"
case "$LOCK_PID" in
  ''|*[!0-9]*) fail "LOCK_PID is not a plain number: $LOCK_PID" ;;
esac
[ "$LOCK_PID" = "$actual_pid" ] || fail "LOCK_PID ($LOCK_PID) doesn't match the PID actually recorded in $LOCK ($actual_pid)"

# PID_RUNNING: must correctly report that this PID isn't alive.
[ -n "${PID_RUNNING:-}" ] || fail "PID_RUNNING is empty"
echo "$PID_RUNNING" | grep -qi '^no$' || fail "PID_RUNNING should be 'no' — PID $actual_pid, the one recorded in the lock, isn't running in this container"

# Ground truth: confirm live that this PID really isn't running. If it
# somehow were, the fixture (or the host) is in a state this lab didn't
# expect.
if kill -0 "$actual_pid" 2>/dev/null; then
  fail "PID $actual_pid is actually alive right now — this lock wouldn't be stale; the check's assumptions may be stale instead"
fi

# LOCK_STATUS: the conclusion that follows from the above.
[ -n "${LOCK_STATUS:-}" ] || fail "LOCK_STATUS is empty"
echo "$LOCK_STATUS" | grep -qi '^stale$' || fail "LOCK_STATUS should be 'stale' — a lock naming a PID that isn't running has nothing left to protect"

echo "PASS: correctly identified deploy.lock (PID $actual_pid, not running) as a stale lock from a dead process"
