#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /workspace || fail "/workspace missing"
git rev-parse --git-dir >/dev/null 2>&1 || fail "/workspace is not a git repository"

[ -d /srv/origin.git ] || fail "/srv/origin.git (the remote) is missing"

ORIGIN_TIP=$(git --git-dir=/srv/origin.git rev-parse main 2>/dev/null) || fail "origin's main ref is missing"
LOCAL_TIP=$(git rev-parse main 2>/dev/null) || fail "local main ref is missing"

[ "$LOCAL_TIP" = "$ORIGIN_TIP" ] || fail "local main and origin/main don't match — push your commit back to origin"

git --git-dir=/srv/origin.git log main --oneline 2>/dev/null | grep -q "Add logging helper" || \
  fail "the teammate's 'Add logging helper' commit is missing from origin — don't lose it, bring it in and push on top of it"

COMMIT_COUNT=$(git --git-dir=/srv/origin.git rev-list --count main)
[ "$COMMIT_COUNT" -ge 4 ] || fail "origin/main should have your own new commit on top of the teammate's, not just the teammate's alone"

echo "PASS: local main and origin/main are in sync, including the teammate's commit and your own"
