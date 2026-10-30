#!/usr/bin/env bash
set -uo pipefail

cd /workspace || exit 1
ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

actual_branch=$(git rev-parse --abbrev-ref HEAD)
[ -n "${CURRENT_BRANCH:-}" ] || fail "CURRENT_BRANCH is empty"
[ "$CURRENT_BRANCH" = "$actual_branch" ] || fail "CURRENT_BRANCH doesn't match the checked-out branch ($actual_branch)"

actual_sha=$(git rev-parse --short HEAD)
[ -n "${HEAD_COMMIT_SHORT_SHA:-}" ] || fail "HEAD_COMMIT_SHORT_SHA is empty"
[ "$HEAD_COMMIT_SHORT_SHA" = "$actual_sha" ] || fail "HEAD_COMMIT_SHORT_SHA doesn't match git rev-parse --short HEAD ($actual_sha)"

[ -n "${OTHER_BRANCH_NAME:-}" ] || fail "OTHER_BRANCH_NAME is empty"
other_branches=$(git branch --format='%(refname:short)' | grep -v -x "$actual_branch")
echo "$other_branches" | grep -qx "$OTHER_BRANCH_NAME" || fail "OTHER_BRANCH_NAME doesn't name a branch other than the checked-out one"

actual_modified=$(git status --porcelain | grep -c '^ M ')
[ -n "${MODIFIED_FILE_COUNT:-}" ] || fail "MODIFIED_FILE_COUNT is empty"
[ "$MODIFIED_FILE_COUNT" = "$actual_modified" ] || fail "MODIFIED_FILE_COUNT doesn't match the number of modified tracked files ($actual_modified)"

echo "PASS: recorded repo state matches the live repository"
