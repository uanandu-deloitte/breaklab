#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

live_reason=$(kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}' 2>/dev/null)
live_exit=$(kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.exitCode}' 2>/dev/null)
live_limit=$(kubectl get pod mem-demo -o jsonpath='{.spec.containers[0].resources.limits.memory}' 2>/dev/null)

[ "$live_reason" = "OOMKilled" ] || fail "mem-demo's last termination reason is '${live_reason:-unknown}', not OOMKilled — the lab's assumptions may be stale"

[ -n "${TERMINATION_REASON:-}" ] || fail "TERMINATION_REASON is empty"
[ "$TERMINATION_REASON" = "$live_reason" ] || fail "TERMINATION_REASON ($TERMINATION_REASON) doesn't match the live termination reason ($live_reason)"

[ -n "${EXIT_CODE:-}" ] || fail "EXIT_CODE is empty"
[ "$EXIT_CODE" = "$live_exit" ] || fail "EXIT_CODE ($EXIT_CODE) doesn't match the live exit code ($live_exit)"

[ -n "${MEMORY_LIMIT:-}" ] || fail "MEMORY_LIMIT is empty"
[ "$MEMORY_LIMIT" = "$live_limit" ] || fail "MEMORY_LIMIT ($MEMORY_LIMIT) doesn't match the container's live configured limit ($live_limit)"

[ -n "${IS_OOM:-}" ] || fail "IS_OOM is empty"
echo "$IS_OOM" | grep -qi '^yes$' || fail "IS_OOM should be 'yes' — the reason field reads OOMKilled, the kernel's own OOM killer, not a plain application crash"

echo "PASS: mem-demo was correctly diagnosed as OOMKilled (exit $live_exit) against a $live_limit memory limit"
