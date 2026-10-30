#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

cd /workspace || fail "/workspace is missing"
[ -d .git ] || fail "/workspace is not a git repository"

git rev-parse v1.0.1 >/dev/null 2>&1 || fail "tag v1.0.1 doesn't exist"
tag_type=$(git cat-file -t v1.0.1 2>/dev/null)
[ "$tag_type" = "tag" ] || fail "v1.0.1 isn't an annotated tag (git cat-file -t v1.0.1 returned '$tag_type', expected 'tag')"

target=$(git rev-list -n1 v1.0.1)

# The hotfix content must be present in the tagged commit's tree.
git show "$target:handler.py" 2>/dev/null | grep -q "HOTFIX-NULLCHECK-MARKER" \
  || fail "the commit tagged v1.0.1 doesn't contain the null-pointer-crash hotfix in handler.py"

# Only the hotfix should have come along -- not the other main-only commits.
if git show "$target:ratelimit.py" >/dev/null 2>&1; then
  fail "the commit tagged v1.0.1 includes ratelimit.py -- the unrelated rate-limiting commit shouldn't have come along"
fi
git show "$target:CHANGELOG.md" 2>/dev/null | grep -q "CHANGELOG-ONLY-MARKER" \
  && fail "the commit tagged v1.0.1 includes the unrelated changelog update -- only the hotfix should have been cherry-picked"

# v1.0.1 should actually be new history, not a re-tag of the old release.
old_target=$(git rev-list -n1 v1.0.0 2>/dev/null) || fail "tag v1.0.0 (the original release) is missing"
[ "$target" != "$old_target" ] || fail "v1.0.1 points at the same commit as v1.0.0 -- the hotfix was never actually cherry-picked"

echo "PASS: v1.0.1 is an annotated tag on a commit that has the hotfix and nothing else from main"
