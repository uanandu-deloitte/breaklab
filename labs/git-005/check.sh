#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /workspace || fail "/workspace missing"
git rev-parse --git-dir >/dev/null 2>&1 || fail "/workspace is not a git repository"

ORIGINAL_TIP=$(git rev-parse original-main-tip 2>/dev/null) || fail "original-main-tip tag missing (fixture corrupted)"
KNOWN_GOOD=$(git rev-parse known-good-base 2>/dev/null) || fail "known-good-base tag missing (fixture corrupted)"

MAIN_TIP=$(git rev-parse main 2>/dev/null) || fail "main branch not found"
[ "$MAIN_TIP" = "$KNOWN_GOOD" ] || fail "main should be back at known-good-base — it currently points somewhere else"

FOUND=""
while read -r branch; do
  [ "$branch" = "main" ] && continue
  tip=$(git rev-parse "refs/heads/$branch" 2>/dev/null) || continue
  if [ "$tip" = "$ORIGINAL_TIP" ]; then
    FOUND="$branch"
    break
  fi
done < <(git for-each-ref --format='%(refname:short)' refs/heads/)

[ -n "$FOUND" ] || fail "no branch (other than main) has a tip matching the original two commits — create one before moving main"

echo "PASS: branch '$FOUND' preserves the original commits and main is back at known-good-base"
