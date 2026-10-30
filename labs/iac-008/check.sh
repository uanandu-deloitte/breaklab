#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
RESDIR=/infra/resources
ORIG_STATE=/infra/.original-state.json
STATE=/infra/state.json
fail() { echo "FAIL: $1"; exit 1; }

[ -d "$RESDIR" ] || fail "/infra/resources is missing — fixture was tampered with"
[ -f "$ORIG_STATE" ] || fail "/infra/.original-state.json is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

# Find which declared path was originally missing from state.json's
# managed list — ground truth, computed from the frozen pre-fix state.
unmanaged=""
for f in "$RESDIR"/*.yaml; do
  p=$(sed -n 's/^path: *//p' "$f")
  if ! jq -e --arg p "$p" '.managed | index($p) != null' "$ORIG_STATE" >/dev/null; then
    unmanaged="$p"
  fi
done
[ -n "$unmanaged" ] || fail "fixture is broken: every declared resource was already managed from the start"

[ -n "${UNMANAGED_PATH:-}" ] || fail "UNMANAGED_PATH is empty"
[ "$UNMANAGED_PATH" = "$unmanaged" ] || fail "UNMANAGED_PATH should be the declared resource that was never in state ($unmanaged), got: $UNMANAGED_PATH"

[ -n "${APPLY_SUCCEEDED_AFTER_IMPORT:-}" ] || fail "APPLY_SUCCEEDED_AFTER_IMPORT is empty"
echo "$APPLY_SUCCEEDED_AFTER_IMPORT" | grep -qi '^yes$' || fail "APPLY_SUCCEEDED_AFTER_IMPORT should be 'yes' — after importing, infra-apply should have actually written the resource"

# Ground truth: the path must now actually be in the live managed list...
jq -e --arg p "$unmanaged" '.managed | index($p) != null' "$STATE" >/dev/null \
  || fail "$unmanaged still isn't in /infra/state.json's managed list — it needs to be imported before infra-apply will ever touch it"

# ...and the real resource must now match its declaration.
decl_file=""
for f in "$RESDIR"/*.yaml; do
  p=$(sed -n 's/^path: *//p' "$f")
  [ "$p" = "$unmanaged" ] && decl_file="$f"
done
[ -n "$decl_file" ] || fail "couldn't find the resource declaration for $unmanaged — fixture was tampered with"
declared_content=$(sed -n 's/^content: *//p' "$decl_file")
[ -f "$unmanaged" ] || fail "real resource $unmanaged is missing"
real_content=$(cat "$unmanaged")
[ "$real_content" = "$declared_content" ] || fail "$unmanaged doesn't match its declared content ($declared_content) — did you re-run infra-apply after importing it?"

echo "PASS: $unmanaged was correctly identified as unmanaged, imported into state, and now matches its declaration"
