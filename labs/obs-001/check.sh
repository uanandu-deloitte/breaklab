#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
[ -f /var/log/app.jsonl ] || fail "/var/log/app.jsonl fixture is missing (don't delete it)"
[ -f /var/lib/metrics/current.json ] || fail "/var/lib/metrics/current.json fixture is missing (don't delete it)"

# shellcheck disable=SC1090
source "$ANSWERS"

expected_error_count=$(jq -r 'select(.level=="ERROR") | .level' /var/log/app.jsonl | wc -l | tr -d ' ')
expected_message=$(jq -r 'select(.request_id=="req-3005") | .message' /var/log/app.jsonl)
expected_disk=$(jq -r '.disk_used_pct' /var/lib/metrics/current.json)

[ -n "${ERROR_COUNT:-}" ] || fail "ERROR_COUNT is empty"
[ "$ERROR_COUNT" = "$expected_error_count" ] || fail "ERROR_COUNT doesn't match the actual number of ERROR lines in app.jsonl (expected $expected_error_count)"

[ -n "${REQUEST_ID_MESSAGE:-}" ] || fail "REQUEST_ID_MESSAGE is empty"
[ "$REQUEST_ID_MESSAGE" = "$expected_message" ] || fail "REQUEST_ID_MESSAGE doesn't match the message on the req-3005 log line"

[ -n "${DISK_USED_PCT:-}" ] || fail "DISK_USED_PCT is empty"
[ "$DISK_USED_PCT" = "$expected_disk" ] || fail "DISK_USED_PCT doesn't match the current disk_used_pct metric value"

[ "${WHY_ANSWER:-}" = "LOG" ] || fail "WHY_ANSWER should be LOG -- a log is the discrete event record that explains why one thing happened"
[ "${HOWMUCH_ANSWER:-}" = "METRIC" ] || fail "HOWMUCH_ANSWER should be METRIC -- a metric is the current numeric value, not an event"

echo "PASS: answers.env correctly distinguishes log events from metric values"
