#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
SRC=/infra/resources.yaml
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$SRC" ] || fail "/infra/resources.yaml is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

file_value=$(sed -n 's/^value: *//p' "$SRC")
env_value="${INFRA_VALUE:-}"
[ -n "$env_value" ] || fail "INFRA_VALUE isn't set in this container's environment — fixture was tampered with"
[ "$file_value" != "$env_value" ] || fail "fixture is broken: file and env value happen to match, so this lab can't test precedence"

[ -n "${FILE_VALUE:-}" ] || fail "FILE_VALUE is empty"
[ "$FILE_VALUE" = "$file_value" ] || fail "FILE_VALUE should be resources.yaml's declared value ($file_value), got: $FILE_VALUE"

[ -n "${ENV_VALUE:-}" ] || fail "ENV_VALUE is empty"
[ "$ENV_VALUE" = "$env_value" ] || fail "ENV_VALUE should be \$INFRA_VALUE's current value ($env_value), got: $ENV_VALUE"

[ -n "${ACTUAL_SOURCE:-}" ] || fail "ACTUAL_SOURCE is empty"
echo "$ACTUAL_SOURCE" | grep -qi '^env$' || fail "ACTUAL_SOURCE should be 'env' — infra-apply prefers the environment variable over the config file whenever no flag is given"

[ -n "${RESOLVED_VALUE:-}" ] || fail "RESOLVED_VALUE is empty"
[ "$RESOLVED_VALUE" = "$env_value" ] || fail "RESOLVED_VALUE should match \$INFRA_VALUE ($env_value) — that's the source infra-apply actually honors here"

# Ground truth: actually run the tool with no flag and confirm it agrees.
live_output=$(infra-apply 2>&1)
echo "$live_output" | grep -q "(source: env):" || fail "live infra-apply run didn't report using the env source — got: $live_output"
echo "$live_output" | grep -qF "$env_value" || fail "live infra-apply run didn't resolve to \$INFRA_VALUE's value — got: $live_output"

path=$(sed -n 's/^path: *//p' "$SRC")
real=$(cat "$path" 2>/dev/null || true)
[ "$real" = "$env_value" ] || fail "$path doesn't currently hold \$INFRA_VALUE's value ($env_value) — got: $real"

echo "PASS: file and environment both declared a value, and answers.env correctly identifies the environment variable as the one infra-apply actually honors"
