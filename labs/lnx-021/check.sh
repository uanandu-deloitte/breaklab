#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${LARGEST_DIR:-}" ] || fail "LARGEST_DIR is empty"
[ -n "${LARGEST_FILE:-}" ] || fail "LARGEST_FILE is empty"

actual_dir=$(du -sb /data/*/ 2>/dev/null | sort -rn | head -1 | awk '{print $2}')
actual_dir=${actual_dir%/}

actual_file=$(find /data -type f -printf '%s %p\n' | sort -rn | head -1 | awk '{print $2}')

given_dir=${LARGEST_DIR%/}

if [ "$given_dir" != "$actual_dir" ]; then
  fail "LARGEST_DIR=$LARGEST_DIR doesn't match the actual largest top-level directory ($actual_dir)"
fi

if [ "$LARGEST_FILE" != "$actual_file" ]; then
  fail "LARGEST_FILE=$LARGEST_FILE doesn't match the actual largest file ($actual_file)"
fi

echo "PASS: LARGEST_DIR ($actual_dir) and LARGEST_FILE ($actual_file) match the real filesystem"
