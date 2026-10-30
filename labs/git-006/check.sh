#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /workspace || fail "/workspace missing"
git rev-parse --git-dir >/dev/null 2>&1 || fail "/workspace is not a git repository"

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
[ "$CURRENT_BRANCH" = "main" ] || fail "not currently on main — check out main after merging"

[ -f signup.html ] || fail "signup.html is missing from main — merge feature/signup-form in"
[ -f validate.js ] || fail "validate.js is missing from main — merge feature/signup-form in"
grep -qi "setup instructions" README.md 2>/dev/null || fail "README.md should still contain main's own setup-instructions commit"

MERGE_COUNT=$(git rev-list --min-parents=2 --count main)
[ "$MERGE_COUNT" -ge 1 ] || fail "main has no merge commit — a fast-forward wouldn't produce one; this merge should"

echo "PASS: feature/signup-form is merged into main with a true merge commit"
