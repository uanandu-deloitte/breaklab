#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
ARTIFACTS_DIR=/srv/artifacts
CURRENT_LINK=$ARTIFACTS_DIR/current
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
[ -L "$CURRENT_LINK" ] || fail "$CURRENT_LINK is missing or not a symlink — fixture was tampered with"

# shellcheck disable=SC1090
source "$ANSWERS"

actual_target="$(readlink "$CURRENT_LINK")"
[ -n "$actual_target" ] || fail "couldn't read the symlink target of $CURRENT_LINK"

# PROMOTED_ARTIFACT: must match what the symlink actually points to right
# now, read live rather than assumed from filenames.
[ -n "${PROMOTED_ARTIFACT:-}" ] || fail "PROMOTED_ARTIFACT is empty"
[ "$PROMOTED_ARTIFACT" = "$actual_target" ] \
  || fail "PROMOTED_ARTIFACT ($PROMOTED_ARTIFACT) doesn't match what $CURRENT_LINK actually points to ($actual_target)"

# Ground truth: the symlink must NOT point at the newest-by-version
# filename, or this lab isn't testing anything.
newest_by_name="$(cd "$ARTIFACTS_DIR" && ls app-*.tar 2>/dev/null | sort -V | tail -n1)"
[ -n "$newest_by_name" ] || fail "couldn't find any app-*.tar artifacts — fixture was tampered with"
[ "$newest_by_name" != "$actual_target" ] \
  || fail "fixture problem: the symlink happens to point at the newest-looking filename ($newest_by_name), so this lab isn't testing anything"

# WOULD_GUESSING_BY_NAME_BE_RIGHT: the whole point — the newest-looking
# filename is not what's actually promoted.
[ -n "${WOULD_GUESSING_BY_NAME_BE_RIGHT:-}" ] || fail "WOULD_GUESSING_BY_NAME_BE_RIGHT is empty"
echo "$WOULD_GUESSING_BY_NAME_BE_RIGHT" | grep -qi '^no$' \
  || fail "WOULD_GUESSING_BY_NAME_BE_RIGHT should be 'no' — the newest-looking filename ($newest_by_name) is not what's actually marked for promotion"

echo "PASS: correctly identified $actual_target as what $CURRENT_LINK actually promotes, not the newest-looking filename ($newest_by_name)"
