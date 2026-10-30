#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
[ -f /var/log/legacy.log ] || fail "/var/log/legacy.log fixture is missing (don't delete it)"
[ -f /var/log/app.jsonl ] || fail "/var/log/app.jsonl fixture is missing (don't delete it)"

# shellcheck disable=SC1090
source "$ANSWERS"

expected_attempts_legacy=$(grep -c "ERROR could not connect to db" /var/log/legacy.log)
expected_final_ts_legacy=$(grep "checkout failed for order ord-5521" /var/log/legacy.log | awk '{print $1, $2}')

expected_attempts_json=$(jq -r 'select(.order_id=="ord-5521" and .level=="ERROR" and (.message|startswith("could not connect to db")))|.timestamp' /var/log/app.jsonl | wc -l | tr -d ' ')
expected_final_ts_json=$(jq -r 'select(.order_id=="ord-5521" and (.message|startswith("checkout failed")))|.timestamp' /var/log/app.jsonl)

[ -n "${ATTEMPTS_FROM_LEGACY:-}" ] || fail "ATTEMPTS_FROM_LEGACY is empty"
[ "$ATTEMPTS_FROM_LEGACY" = "$expected_attempts_legacy" ] || fail "ATTEMPTS_FROM_LEGACY doesn't match legacy.log (expected $expected_attempts_legacy)"

[ -n "${FINAL_FAILURE_TIME_FROM_LEGACY:-}" ] || fail "FINAL_FAILURE_TIME_FROM_LEGACY is empty"
[ "$FINAL_FAILURE_TIME_FROM_LEGACY" = "$expected_final_ts_legacy" ] || fail "FINAL_FAILURE_TIME_FROM_LEGACY doesn't match legacy.log's final failure line (expected '$expected_final_ts_legacy')"

[ -n "${METHOD_LEGACY:-}" ] || fail "METHOD_LEGACY is empty"
if echo "$METHOD_LEGACY" | grep -qi "jq"; then
  fail "METHOD_LEGACY should name a text tool (grep/awk/sed), not jq -- that's the structured-log one"
fi

[ -n "${ATTEMPTS_FROM_JSON:-}" ] || fail "ATTEMPTS_FROM_JSON is empty"
[ "$ATTEMPTS_FROM_JSON" = "$expected_attempts_json" ] || fail "ATTEMPTS_FROM_JSON doesn't match app.jsonl (expected $expected_attempts_json)"

[ -n "${FINAL_FAILURE_TIME_FROM_JSON:-}" ] || fail "FINAL_FAILURE_TIME_FROM_JSON is empty"
[ "$FINAL_FAILURE_TIME_FROM_JSON" = "$expected_final_ts_json" ] || fail "FINAL_FAILURE_TIME_FROM_JSON doesn't match app.jsonl's final failure line (expected '$expected_final_ts_json')"

[ -n "${METHOD_JSON:-}" ] || fail "METHOD_JSON is empty"
if ! echo "$METHOD_JSON" | grep -qi "jq"; then
  fail "METHOD_JSON should mention jq -- that's the tool for querying structured app.jsonl"
fi

echo "PASS: both the unstructured and structured extraction agree on the real failure"
