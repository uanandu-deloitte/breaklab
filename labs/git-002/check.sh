#!/usr/bin/env bash
set -uo pipefail

cd /workspace || exit 1
fail() { echo "FAIL: $1"; exit 1; }

head_msg=$(git log -1 --pretty=%s 2>/dev/null) || fail "no commits found on top of the fixture"
[ "$head_msg" = "Increase API timeout to 30s" ] || fail "HEAD commit message is not exactly 'Increase API timeout to 30s' (got: $head_msg)"

changed_files=$(git diff-tree --no-commit-id --name-only -r HEAD)
[ "$changed_files" = "config.yaml" ] || fail "the new commit should touch only config.yaml (touched: $changed_files)"

grep -q 'timeout: 30' config.yaml || fail "config.yaml doesn't contain the updated timeout"

if git ls-files --error-unmatch scratch-notes.txt >/dev/null 2>&1; then
  fail "scratch-notes.txt is tracked — it should never have been staged or committed"
fi

[ -f scratch-notes.txt ] || fail "scratch-notes.txt should still exist on disk, just untracked"

status_line=$(git status --porcelain -- scratch-notes.txt)
case "$status_line" in
  '??'*) ;;
  *) fail "scratch-notes.txt should show as untracked (??) in git status, got: $status_line" ;;
esac

echo "PASS: config.yaml was committed alone with the exact message, scratch-notes.txt is still untracked"
