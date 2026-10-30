#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
LOG=/workspace/access.log
ORIG=/opt/fixtures/access.log.original

[ -f "$LOG" ] || fail "/workspace/access.log not found"
if ! diff -q <(sort "$LOG") <(sort "$ORIG") >/dev/null; then
  fail "access.log content has been modified — don't edit the input log"
fi

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

EXPECTED_TOTAL=16
EXPECTED_FIRST="2026-01-10T09:00:01Z"
EXPECTED_LAST="2026-01-10T09:01:15Z"
EXPECTED_2XX=11
EXPECTED_4XX=3
EXPECTED_5XX=2
EXPECTED_PATH="/api/orders"

[ -n "${TOTAL_REQUESTS:-}" ] || fail "TOTAL_REQUESTS is empty"
[ "$TOTAL_REQUESTS" = "$EXPECTED_TOTAL" ] || fail "TOTAL_REQUESTS should be $EXPECTED_TOTAL, got $TOTAL_REQUESTS"

[ -n "${FIRST_TIMESTAMP:-}" ] || fail "FIRST_TIMESTAMP is empty"
[ "$FIRST_TIMESTAMP" = "$EXPECTED_FIRST" ] || fail "FIRST_TIMESTAMP should be $EXPECTED_FIRST, got $FIRST_TIMESTAMP"

[ -n "${LAST_TIMESTAMP:-}" ] || fail "LAST_TIMESTAMP is empty"
[ "$LAST_TIMESTAMP" = "$EXPECTED_LAST" ] || fail "LAST_TIMESTAMP should be $EXPECTED_LAST, got $LAST_TIMESTAMP"

[ -n "${COUNT_2XX:-}" ] || fail "COUNT_2XX is empty"
[ "$COUNT_2XX" = "$EXPECTED_2XX" ] || fail "COUNT_2XX should be $EXPECTED_2XX, got $COUNT_2XX"

[ -n "${COUNT_4XX:-}" ] || fail "COUNT_4XX is empty"
[ "$COUNT_4XX" = "$EXPECTED_4XX" ] || fail "COUNT_4XX should be $EXPECTED_4XX, got $COUNT_4XX"

[ -n "${COUNT_5XX:-}" ] || fail "COUNT_5XX is empty"
[ "$COUNT_5XX" = "$EXPECTED_5XX" ] || fail "COUNT_5XX should be $EXPECTED_5XX, got $COUNT_5XX"

[ -n "${MOST_FREQUENT_PATH:-}" ] || fail "MOST_FREQUENT_PATH is empty"
[ "$MOST_FREQUENT_PATH" = "$EXPECTED_PATH" ] || fail "MOST_FREQUENT_PATH should be $EXPECTED_PATH, got $MOST_FREQUENT_PATH"

echo "PASS: answers.env matches the access log exactly"
