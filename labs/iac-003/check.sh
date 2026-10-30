#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
SRC=/infra/resources.yaml
ORIG=/infra/.original-value

fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -f "$SRC" ] || fail "/infra/resources.yaml is missing"
[ -f "$ORIG" ] || fail "/infra/.original-value snapshot is missing"

path=$(sed -n 's/^path: *//p' "$SRC")
declared_value=$(sed -n 's/^value: *//p' "$SRC")
original_value=$(cat "$ORIG")

[ -n "$path" ] || fail "resources.yaml has no path: line"
[ -n "$declared_value" ] || fail "resources.yaml has no value: line"

[ -n "${FIRST_RUN_RESULT:-}" ] || fail "FIRST_RUN_RESULT is empty"
[ -n "${SECOND_RUN_RESULT:-}" ] || fail "SECOND_RUN_RESULT is empty"
[ -n "${FINAL_VALUE:-}" ] || fail "FINAL_VALUE is empty"

echo "$FIRST_RUN_RESULT" | grep -qiE "no.?change|nothing|up.?to.?date" \
  || fail "FIRST_RUN_RESULT should describe the first run as a no-op (got: $FIRST_RUN_RESULT)"

echo "$SECOND_RUN_RESULT" | grep -qiE "no.?change|nothing|up.?to.?date" \
  && fail "SECOND_RUN_RESULT should NOT describe a no-op — the edit was supposed to force a real change (got: $SECOND_RUN_RESULT)"

echo "$SECOND_RUN_RESULT" | grep -qiE "^changed$|change" \
  || fail "SECOND_RUN_RESULT should describe a real change (got: $SECOND_RUN_RESULT)"

[ "$declared_value" != "$original_value" ] \
  || fail "resources.yaml's value still matches the original build-time value — it was never actually edited"

[ "$FINAL_VALUE" = "$declared_value" ] \
  || fail "FINAL_VALUE ($FINAL_VALUE) doesn't match resources.yaml's current declared value ($declared_value)"

[ -f "$path" ] || fail "resource file $path does not exist"
actual=$(cat "$path")
[ "$actual" = "$declared_value" ] \
  || fail "resource $path holds '$actual', which doesn't match resources.yaml's current value '$declared_value' — did infra-apply run after your edit?"

echo "PASS: infra-apply was a no-op first, a real change second, and $path now matches resources.yaml"
