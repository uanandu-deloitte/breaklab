#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${HIDDEN_FILES_OPTION:-}" ] || fail "HIDDEN_FILES_OPTION is empty"
echo "$HIDDEN_FILES_OPTION" | grep -qi 'a' || fail "HIDDEN_FILES_OPTION doesn't look like ls's hidden-files flag"

[ -n "${NESTED_MKDIR_OPTION:-}" ] || fail "NESTED_MKDIR_OPTION is empty"
echo "$NESTED_MKDIR_OPTION" | grep -qi 'p' || fail "NESTED_MKDIR_OPTION doesn't look like mkdir's parents flag"

[ -n "${INTERACTIVE_COPY_OPTION:-}" ] || fail "INTERACTIVE_COPY_OPTION is empty"
echo "$INTERACTIVE_COPY_OPTION" | grep -qi 'i' || fail "INTERACTIVE_COPY_OPTION doesn't look like cp's interactive flag"

[ -f /workspace/.hidden-check ] || fail "/workspace/.hidden-check (a hidden file) is missing"

[ -d /workspace/deep/nested/tree ] || fail "/workspace/deep/nested/tree directory tree is missing"

[ -f /workspace/fixture-copy.txt ] || fail "/workspace/fixture-copy.txt is missing"
[ -f /opt/fixtures/readme.txt ] || fail "/opt/fixtures/readme.txt fixture is missing (don't delete it)"
cmp -s /workspace/fixture-copy.txt /opt/fixtures/readme.txt || fail "/workspace/fixture-copy.txt doesn't match the original fixture content"

echo "PASS: discovered options are recorded and the operations were performed"
