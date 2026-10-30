#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
SRC=/infra/resources.yaml
ORIG_LOCK=/infra/.original-lock
LOCK=/var/run/infra/apply.lock
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$SRC" ] || fail "/infra/resources.yaml is missing — fixture was tampered with"
[ -f "$ORIG_LOCK" ] || fail "/infra/.original-lock is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

orig_pid=$(sed -n 's/^PID=//p' "$ORIG_LOCK")

[ -n "${LOCK_PID:-}" ] || fail "LOCK_PID is empty"
[ "$LOCK_PID" = "$orig_pid" ] || fail "LOCK_PID should be the PID recorded in the original stale lock ($orig_pid), got: $LOCK_PID"

[ -n "${LOCK_PROCESS_ALIVE:-}" ] || fail "LOCK_PROCESS_ALIVE is empty"
echo "$LOCK_PROCESS_ALIVE" | grep -qi '^no$' || fail "LOCK_PROCESS_ALIVE should be 'no' — PID $orig_pid belongs to no process that could possibly be running in this container"
ps -p "$orig_pid" >/dev/null 2>&1 && fail "ground truth check itself found PID $orig_pid alive — lab fixture assumptions are stale"

[ -n "${APPLY_SUCCEEDED:-}" ] || fail "APPLY_SUCCEEDED is empty"
echo "$APPLY_SUCCEEDED" | grep -qi '^yes$' || fail "APPLY_SUCCEEDED should be 'yes' — infra-apply should have run successfully once the stale lock was cleared"

[ ! -f "$LOCK" ] || fail "$LOCK still exists — a lock left behind after a successful run just becomes the next stale lock"

path=$(sed -n 's/^path: *//p' "$SRC")
declared=$(sed -n 's/^content: *//p' "$SRC")
[ -f "$path" ] || fail "real resource $path is missing"
real=$(cat "$path")
[ "$real" = "$declared" ] || fail "$path doesn't match resources.yaml's declared content — infra-apply hasn't actually succeeded yet"

echo "PASS: stale lock correctly identified and cleared, and infra-apply has actually applied the real resource"
