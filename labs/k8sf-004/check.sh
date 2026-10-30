#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

pod=$(kubectl get pods -l app=flaky-worker -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
[ -n "$pod" ] || fail "flaky-worker Pod not found"

live_restarts=$(kubectl get pod "$pod" -o jsonpath='{.status.containerStatuses[0].restartCount}' 2>/dev/null)
[ -n "$live_restarts" ] || fail "could not read restartCount for $pod"

[ -n "${RESTART_COUNT:-}" ] || fail "RESTART_COUNT is empty"
[ "$RESTART_COUNT" = "$live_restarts" ] || fail "RESTART_COUNT ($RESTART_COUNT) doesn't match the Pod's live restartCount ($live_restarts)"

live_reason=$(kubectl get pod "$pod" -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}' 2>/dev/null)
[ -n "$live_reason" ] || fail "could not read lastState.terminated.reason for $pod"

[ -n "${PREVIOUS_EXIT_REASON:-}" ] || fail "PREVIOUS_EXIT_REASON is empty"
[ "$PREVIOUS_EXIT_REASON" = "$live_reason" ] || fail "PREVIOUS_EXIT_REASON ($PREVIOUS_EXIT_REASON) doesn't match the Pod's live lastState.terminated.reason ($live_reason)"

echo "PASS: recorded restart count ($live_restarts) and previous exit reason ($live_reason) match the live Pod ($pod)"
