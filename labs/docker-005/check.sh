#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

EXPECTED_CONTENT="Welcome to the shared data volume."

# Verify the file's actual content on disk directly — don't just trust the
# recorded yes/no answer below.
[ -f /data/welcome.txt ] || fail "/data/welcome.txt is missing — the image's own copy is shadowed by the mounted volume, so it must be recreated at runtime"
actual_content=$(cat /data/welcome.txt)
[ "$actual_content" = "$EXPECTED_CONTENT" ] || fail "/data/welcome.txt content doesn't match the original (got: $actual_content)"

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${WELCOME_FILE_RESTORED:-}" ] || fail "WELCOME_FILE_RESTORED is empty"
echo "$WELCOME_FILE_RESTORED" | grep -qi '^yes$' || fail "WELCOME_FILE_RESTORED should be yes — the file is restored above"

# Compute the real mount comparison dynamically, never hardcoded, so this
# check stays correct across Docker backends.
data_dev=$(stat -c '%d' /data)
root_dev=$(stat -c '%d' /)
real_distinct=no
[ "$data_dev" != "$root_dev" ] && real_distinct=yes

[ -n "${MOUNT_IS_DISTINCT:-}" ] || fail "MOUNT_IS_DISTINCT is empty"
answer_norm=$(echo "$MOUNT_IS_DISTINCT" | tr '[:upper:]' '[:lower:]')
[ "$answer_norm" = "$real_distinct" ] || fail "MOUNT_IS_DISTINCT ($MOUNT_IS_DISTINCT) doesn't match the actual stat comparison (data dev=$data_dev, root dev=$root_dev -> $real_distinct)"

echo "PASS: /data/welcome.txt is genuinely restored and the mount comparison is correct"
