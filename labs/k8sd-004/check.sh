#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${DEPENDENCY_NAME:-}" ] || fail "DEPENDENCY_NAME is empty"
logged_dep=$(kubectl logs deployment/frontend --tail=20 2>/dev/null | grep -oE '^[a-zA-Z0-9_-]+: unreachable$' | tail -n1 | sed 's/: unreachable$//')
[ -n "$logged_dep" ] || fail "frontend's logs don't currently report any dependency as unreachable — is the lab actually broken?"
echo "$DEPENDENCY_NAME" | grep -qi "^${logged_dep}$" || fail "DEPENDENCY_NAME ($DEPENDENCY_NAME) doesn't match what frontend's own logs report as unreachable ($logged_dep)"

[ -n "${READY_ENDPOINTS:-}" ] || fail "READY_ENDPOINTS is empty"
case "$READY_ENDPOINTS" in
  ''|*[!0-9]*) fail "READY_ENDPOINTS is not a plain number: $READY_ENDPOINTS" ;;
esac
live_endpoints=$(kubectl get endpoints "$logged_dep" -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null | wc -w | tr -d ' ')
[ "$READY_ENDPOINTS" = "$live_endpoints" ] || fail "READY_ENDPOINTS ($READY_ENDPOINTS) doesn't match kubectl get endpoints $logged_dep ($live_endpoints ready addresses)"

[ -n "${DESIRED_REPLICAS:-}" ] || fail "DESIRED_REPLICAS is empty"
case "$DESIRED_REPLICAS" in
  ''|*[!0-9]*) fail "DESIRED_REPLICAS is not a plain number: $DESIRED_REPLICAS" ;;
esac
live_desired=$(kubectl get deployment "$logged_dep" -o jsonpath='{.spec.replicas}' 2>/dev/null)
live_desired=${live_desired:-0}
[ "$DESIRED_REPLICAS" = "$live_desired" ] || fail "DESIRED_REPLICAS ($DESIRED_REPLICAS) doesn't match kubectl get deployment $logged_dep's DESIRED ($live_desired)"

[ -n "${READY_REPLICAS:-}" ] || fail "READY_REPLICAS is empty"
case "$READY_REPLICAS" in
  ''|*[!0-9]*) fail "READY_REPLICAS is not a plain number: $READY_REPLICAS" ;;
esac
live_ready=$(kubectl get deployment "$logged_dep" -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
live_ready=${live_ready:-0}
[ "$READY_REPLICAS" = "$live_ready" ] || fail "READY_REPLICAS ($READY_REPLICAS) doesn't match kubectl get deployment $logged_dep's READY ($live_ready)"

[ -n "${FAILURE_MODE:-}" ] || fail "FAILURE_MODE is empty"
if [ "$live_desired" = "0" ]; then
  expected_mode="scaled-to-zero"
else
  expected_mode="crash-looping"
fi
echo "$FAILURE_MODE" | grep -qi "^${expected_mode}$" || fail "FAILURE_MODE ($FAILURE_MODE) should be '$expected_mode' given DESIRED=$live_desired and READY=$live_ready"

echo "PASS: correctly diagnosed $logged_dep as $expected_mode (desired=$live_desired, ready=$live_ready, endpoints=$live_endpoints)"
