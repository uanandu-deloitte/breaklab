#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

ORIGIN=/srv/origin.git
WS=/workspace

[ -d "$ORIGIN" ] || fail "$ORIGIN (the bare origin repo) is missing"
[ -d "$WS/.git" ] || fail "$WS is not a git repository"

TEAMMATE_SHA=$(cat /srv/.teammate-sha 2>/dev/null) || fail "missing fixture record of the teammate's commit"
LOCAL_SHA=$(cat /srv/.local-sha 2>/dev/null) || fail "missing fixture record of your original local commit"

origin_tip=$(git --git-dir="$ORIGIN" rev-parse main 2>/dev/null) || fail "origin has no main branch"
local_tip=$(git -C "$WS" rev-parse main 2>/dev/null) || fail "local /workspace has no main branch"

[ "$origin_tip" = "$local_tip" ] || fail "local main ($local_tip) doesn't match origin's main ($origin_tip) -- push a fast-forward so both sides agree, don't leave them diverged"

git --git-dir="$ORIGIN" merge-base --is-ancestor "$TEAMMATE_SHA" "$origin_tip" \
  || fail "origin's main no longer contains the teammate's commit -- don't force-push over it"

if ! git --git-dir="$ORIGIN" merge-base --is-ancestor "$LOCAL_SHA" "$origin_tip" 2>/dev/null; then
  # A rebase rewrites the local commit's SHA, so fall back to checking its
  # content landed instead of insisting on the exact original commit object.
  git --git-dir="$ORIGIN" show "$origin_tip:app.py" 2>/dev/null | grep -q "LOCAL-COMMIT-MARKER" \
    || fail "origin's main is missing your local commit's change (app.py has no LOCAL-COMMIT-MARKER)"
fi

echo "PASS: origin's main contains both the teammate's commit and your local commit, and local main matches it exactly"
