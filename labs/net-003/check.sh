#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

[ "${ENDPOINT_1:-}" = "healthy" ]     || fail "ENDPOINT_1 (port 8001) should be classified 'healthy'"
[ "${ENDPOINT_2:-}" = "dns-failure" ] || fail "ENDPOINT_2 (nonexistent.internal:8002) should be classified 'dns-failure'"
[ "${ENDPOINT_3:-}" = "refused" ]     || fail "ENDPOINT_3 (port 8003) should be classified 'refused'"
[ "${ENDPOINT_4:-}" = "timeout" ]     || fail "ENDPOINT_4 (port 8004) should be classified 'timeout'"

echo "PASS: all four endpoints correctly classified"
