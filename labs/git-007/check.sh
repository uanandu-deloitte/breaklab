#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /workspace || fail "/workspace missing"
git rev-parse --git-dir >/dev/null 2>&1 || fail "/workspace is not a git repository"

if grep -rl '<<<<<<<' /workspace >/dev/null 2>&1; then
  fail "conflict markers (<<<<<<<) are still present somewhere under /workspace — resolve them before committing"
fi

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
[ "$CURRENT_BRANCH" = "main" ] || fail "not currently on main — complete the merge from main"

MERGE_COUNT=$(git rev-list --min-parents=2 --count main)
[ "$MERGE_COUNT" -ge 1 ] || fail "main has no merge commit yet — finish the merge with git commit after resolving"

[ -f config.yaml ] || fail "config.yaml is missing"
EXPECTED=$'service: acme-api\nport: 8080\ntimeout: 60'
ACTUAL=$(cat config.yaml)
[ "$ACTUAL" = "$EXPECTED" ] || fail "config.yaml doesn't match the agreed resolution (service: acme-api / port: 8080 / timeout: 60)"

echo "PASS: conflict resolved to the agreed content and the merge commit is in place"
