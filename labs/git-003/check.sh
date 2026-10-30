#!/usr/bin/env bash
set -uo pipefail

cd /workspace || exit 1
fail() { echo "FAIL: $1"; exit 1; }

git rev-parse fixture-base >/dev/null 2>&1 || fail "fixture-base tag missing (fixture is broken)"

[ -z "$(git status --porcelain)" ] || fail "working tree has uncommitted changes — commit everything first"

commit_count=$(git rev-list fixture-base..HEAD --count)
[ "$commit_count" -ge 2 ] || fail "expected at least 2 new commits since fixture-base, found $commit_count"

grep -q 'len(items) - 1' app.py || fail "app.py's get_last_item still has the off-by-one bug"
grep -q 'DEBUG: utils.py loaded' utils.py || fail "utils.py should still contain the debug print, committed in its own commit"

mixed=""
app_commit_found=0
utils_commit_found=0
for sha in $(git rev-list fixture-base..HEAD); do
  files=$(git diff-tree --no-commit-id --name-only -r "$sha")
  touches_app=0
  touches_utils=0
  echo "$files" | grep -qx "app.py" && touches_app=1 && app_commit_found=1
  echo "$files" | grep -qx "utils.py" && touches_utils=1 && utils_commit_found=1
  if [ "$touches_app" = 1 ] && [ "$touches_utils" = 1 ]; then
    mixed="$sha"
  fi
done

[ -n "$mixed" ] && fail "commit $mixed touches both app.py and utils.py — split the fix and the unrelated change into separate commits"
[ "$app_commit_found" = 1 ] || fail "no commit touches app.py"
[ "$utils_commit_found" = 1 ] || fail "no commit touches utils.py"

echo "PASS: the bugfix and the unrelated change are in separate commits"
