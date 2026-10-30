#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
SRC=/infra/resources.yaml
ORIG_PIN=/infra/.original-pin
PIN=/infra/MODULE_VERSION
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$SRC" ] || fail "/infra/resources.yaml is missing — fixture was tampered with"
[ -f "$ORIG_PIN" ] || fail "/infra/.original-pin is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

declared=$(sed -n 's/^module: *//p' "$SRC")
original_pin=$(cat "$ORIG_PIN")

[ -n "${DECLARED_MODULE:-}" ] || fail "DECLARED_MODULE is empty"
[ "$DECLARED_MODULE" = "$declared" ] || fail "DECLARED_MODULE should be resources.yaml's declared module ($declared), got: $DECLARED_MODULE"

[ -n "${PINNED_MODULE:-}" ] || fail "PINNED_MODULE is empty"
[ "$PINNED_MODULE" = "$original_pin" ] || fail "PINNED_MODULE should be the module MODULE_VERSION was originally pinned to ($original_pin), got: $PINNED_MODULE"

[ -n "${MODULES_AGREED_BEFORE_FIX:-}" ] || fail "MODULES_AGREED_BEFORE_FIX is empty"
if [ "$declared" = "$original_pin" ]; then
  echo "$MODULES_AGREED_BEFORE_FIX" | grep -qi '^yes$' || fail "MODULES_AGREED_BEFORE_FIX should be 'yes' — the declared module and the original pin actually matched"
else
  echo "$MODULES_AGREED_BEFORE_FIX" | grep -qi '^no$' || fail "MODULES_AGREED_BEFORE_FIX should be 'no' — the declared module ($declared) and the original pin ($original_pin) disagreed"
fi

[ -n "${RESOLVED_VERSION_AFTER_FIX:-}" ] || fail "RESOLVED_VERSION_AFTER_FIX is empty"
[ "$RESOLVED_VERSION_AFTER_FIX" = "$declared" ] || fail "RESOLVED_VERSION_AFTER_FIX should match the declared module ($declared) once the pin is fixed, got: $RESOLVED_VERSION_AFTER_FIX"

current_pin=$(cat "$PIN" 2>/dev/null || true)
[ "$current_pin" = "$declared" ] || fail "/infra/MODULE_VERSION still says '$current_pin' — fix the pin so it matches the declared module ($declared)"

path=$(sed -n 's/^path: *//p' "$SRC")
[ -f "$path" ] || fail "real resource $path is missing"
real=$(cat "$path")
echo "$real" | grep -q "rendered by v2\]$" || fail "$path doesn't look like it was rendered by $declared — did you re-run infra-apply after fixing the pin?"

echo "PASS: the declared module and the pin now agree, and the deployed resource was actually rendered by $declared"
