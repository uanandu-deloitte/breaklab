#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

SOURCE=/var/log/app.log
DEST=/var/log/aggregator.log

[ -f "$SOURCE" ] || fail "source log $SOURCE is missing"
[ -f "$DEST" ] || fail "destination $DEST does not exist yet — run ship-logs"

diff "$SOURCE" "$DEST" >/dev/null 2>&1 \
  || fail "$DEST doesn't match $SOURCE yet — run ship-logs --source $SOURCE --destination $DEST"

before_count=$(wc -l < "$DEST")

probe="PROBE $(date -u +%FT%TZ) check.sh verifying offset tracking $$"
echo "$probe" >> "$SOURCE"

ship-logs --source "$SOURCE" --destination "$DEST" >/dev/null 2>&1 \
  || fail "ship-logs failed to run with --source/--destination"

after_count=$(wc -l < "$DEST")

if [ "$after_count" -ne "$((before_count + 1))" ]; then
  fail "expected exactly one new line in $DEST after re-running ship-logs, got $((after_count - before_count)) new line(s) — check for duplicates or missing lines"
fi

last_line=$(tail -n 1 "$DEST")
[ "$last_line" = "$probe" ] \
  || fail "the newly shipped line in $DEST doesn't match the probe line appended to $SOURCE"

diff "$SOURCE" "$DEST" >/dev/null 2>&1 \
  || fail "$DEST doesn't exactly match $SOURCE after shipping the probe line"

echo "PASS: ship-logs shipped exactly the new line without duplicating existing ones"
