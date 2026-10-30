#!/usr/bin/env bash
set -uo pipefail

FILE=/srv/team/newfile.txt
fail() { echo "FAIL: $1"; exit 1; }

[ -e "$FILE" ] || fail "$FILE does not exist — create it after setting a safe umask"
[ -f "$FILE" ] || fail "$FILE exists but is not a regular file"

mode=$(stat -c '%a' "$FILE")
other=${mode: -1}
group=${mode: -2:1}

if [ "$other" != "0" ]; then
  fail "$FILE has mode $mode — other still has access (expected no access for other)"
fi

if ! (( group & 4 )); then
  fail "$FILE has mode $mode — group does not have read access"
fi

if ! (( group & 2 )); then
  fail "$FILE has mode $mode — group does not have write access"
fi

echo "PASS: $FILE has mode $mode — group can read/write, other has no access"
