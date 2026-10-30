#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
SRC=/infra/resources.yaml
ORIG=/infra/.original-value
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$SRC" ] || fail "/infra/resources.yaml is missing — fixture was tampered with"
[ -f "$ORIG" ] || fail "/infra/.original-value is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

path=$(sed -n 's/^path: *//p' "$SRC")
declared=$(sed -n 's/^content: *//p' "$SRC")
original=$(cat "$ORIG")

[ -n "$path" ] || fail "couldn't parse path from /infra/resources.yaml"
[ "$declared" != "$original" ] || fail "resources.yaml's content should differ from the original value — fixture was tampered with"

[ -f "$path" ] || fail "real resource $path is missing"
real=$(cat "$path")

[ "$real" = "$declared" ] || fail "$path doesn't match resources.yaml's declared content — did you run infra-apply (with no flag) at the end?"

[ -n "${PLANNED_CHANGE:-}" ] || fail "PLANNED_CHANGE is empty"
echo "$PLANNED_CHANGE" | grep -qF "$original" || fail "PLANNED_CHANGE doesn't mention the original value ('$original') — record infra-apply --plan's actual output, run before the real apply"
echo "$PLANNED_CHANGE" | grep -qF "$declared" || fail "PLANNED_CHANGE doesn't mention the new declared value ('$declared') — record infra-apply --plan's actual output"

echo "PASS: the plan was recorded correctly and $path now matches resources.yaml's declared value"
