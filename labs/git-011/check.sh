#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${STASH_COMMAND:-}" ] || fail "STASH_COMMAND is empty"
echo "$STASH_COMMAND" | grep -qi 'push' || fail "STASH_COMMAND doesn't look like git stash push"

[ -n "${STASH_INCLUDE_UNTRACKED_FLAG:-}" ] || fail "STASH_INCLUDE_UNTRACKED_FLAG is empty"
echo "$STASH_INCLUDE_UNTRACKED_FLAG" | grep -qi -- '-u\|--include-untracked' || fail "STASH_INCLUDE_UNTRACKED_FLAG doesn't look like the untracked flag (-u / --include-untracked)"

[ -n "${POP_COMMAND:-}" ] || fail "POP_COMMAND is empty"
echo "$POP_COMMAND" | grep -qi 'pop' || fail "POP_COMMAND doesn't look like git stash pop"

cd /workspace || fail "/workspace is missing"
[ -d .git ] || fail "/workspace is not a git repository"

[ -f app.py ] || fail "app.py is missing"
grep -q "TRACKED-CHANGE-MARKER" app.py || fail "app.py doesn't have the tracked change restored"

[ -f notes.txt ] || fail "notes.txt (the untracked file) is missing"
grep -q "UNTRACKED-FILE-MARKER" notes.txt || fail "notes.txt doesn't match its original content"

# app.py should show as a modified tracked file (unstaged), not committed.
status_line=$(git status --porcelain -- app.py)
echo "$status_line" | grep -q '^ M' || fail "app.py should show as an unstaged modification (git status), got: ${status_line:-<nothing>}"

# notes.txt should show as untracked.
notes_status=$(git status --porcelain -- notes.txt)
echo "$notes_status" | grep -q '^??' || fail "notes.txt should show as untracked (git status), got: ${notes_status:-<nothing>}"

stash_count=$(git stash list | wc -l | tr -d ' ')
[ "$stash_count" = "0" ] || fail "git stash list still has $stash_count entr(ies) -- pop it, don't leave it on the stack"

echo "PASS: both the tracked change and the untracked file are back, and the stash was popped"
