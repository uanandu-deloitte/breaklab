#!/usr/bin/env bash
set -uo pipefail
ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${FAILING_DEPENDENCY:-}" ] || fail "FAILING_DEPENDENCY is empty"
[ "$FAILING_DEPENDENCY" = "payments" ] || fail "FAILING_DEPENDENCY should name the service whose own upstream call actually timed out, not the one that returned the error to the user"

[ -n "${REQUEST_ID:-}" ] || fail "REQUEST_ID is empty"
[ "$REQUEST_ID" = "req-8841" ] || fail "REQUEST_ID doesn't match the request that failed end to end across all three logs"

echo "PASS: answers.env correctly identifies the failing dependency and request"
