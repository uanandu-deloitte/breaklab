#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
SRC=/infra/resources.yaml
ORIG=/infra/.original-content
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$SRC" ] || fail "/infra/resources.yaml is missing — fixture was tampered with"
[ -f "$ORIG" ] || fail "/infra/.original-content is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

path=$(sed -n 's/^path: *//p' "$SRC")
declared=$(sed -n 's/^content: *//p' "$SRC")
original=$(cat "$ORIG")

[ -n "$path" ] || fail "couldn't parse path from /infra/resources.yaml"
[ "$declared" != "$original" ] || fail "resources.yaml's content hasn't been changed from its original value — edit it first"

[ -f "$path" ] || fail "real resource $path is missing"
real=$(cat "$path")

[ "$real" = "$declared" ] || fail "$path doesn't match resources.yaml's current declared content — did you run infra-apply after editing it?"

[ -n "${SOURCE_OF_TRUTH_FILE:-}" ] || fail "SOURCE_OF_TRUTH_FILE is empty"
[ "$SOURCE_OF_TRUTH_FILE" = "/infra/resources.yaml" ] || fail "SOURCE_OF_TRUTH_FILE should be /infra/resources.yaml, got: $SOURCE_OF_TRUTH_FILE"

[ -n "${DESIRED_VALUE_BEFORE:-}" ] || fail "DESIRED_VALUE_BEFORE is empty"
[ "$DESIRED_VALUE_BEFORE" = "$original" ] || fail "DESIRED_VALUE_BEFORE should be the original declared value ($original), got: $DESIRED_VALUE_BEFORE"

[ -n "${DESIRED_VALUE_AFTER:-}" ] || fail "DESIRED_VALUE_AFTER is empty"
[ "$DESIRED_VALUE_AFTER" = "$declared" ] || fail "DESIRED_VALUE_AFTER should match resources.yaml's current content ($declared), got: $DESIRED_VALUE_AFTER"

echo "PASS: $path matches resources.yaml's current declared value, and answers.env records the change correctly"
