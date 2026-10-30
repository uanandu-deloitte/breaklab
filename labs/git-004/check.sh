#!/usr/bin/env bash
set -uo pipefail

cd /workspace || exit 1
fail() { echo "FAIL: $1"; exit 1; }

git rev-parse original-main-tip >/dev/null 2>&1 || fail "original-main-tip tag missing (fixture is broken)"
original_sha=$(git rev-parse original-main-tip)

git rev-parse --verify feature-logging >/dev/null 2>&1 || fail "branch feature-logging does not exist"

new_commits=$(git rev-list main..feature-logging --count)
[ "$new_commits" -ge 1 ] || fail "feature-logging has no commit that main doesn't have"

git show feature-logging:logging.md >/dev/null 2>&1 || fail "logging.md doesn't exist on feature-logging"

main_sha=$(git rev-parse main)
[ "$main_sha" = "$original_sha" ] || fail "main's tip commit changed — it should still point at the original fixture commit"

current_branch=$(git rev-parse --abbrev-ref HEAD)
[ "$current_branch" = "main" ] || fail "switch back to main when you're done (currently on $current_branch)"

echo "PASS: feature-logging has a new commit and main's tip is unchanged"
