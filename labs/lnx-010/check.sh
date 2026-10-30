#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

SCRIPT=/workspace/top-500-ip.sh
LOG=/workspace/access.log
ORIG=/opt/fixtures/access.log.original
EXPECTED="4 198.51.100.9"

[ -f "$SCRIPT" ] || fail "/workspace/top-500-ip.sh not found"
[ -x "$SCRIPT" ] || fail "/workspace/top-500-ip.sh is not executable (chmod +x it)"
[ -f "$LOG" ] || fail "/workspace/access.log not found"

if ! diff -q <(sort "$LOG") <(sort "$ORIG") >/dev/null; then
  fail "access.log content has been modified — don't edit the input log"
fi

got=$(bash "$SCRIPT" 2>/dev/null | tail -1 | tr -s ' ')
if [ "$got" != "$EXPECTED" ]; then
  fail "expected \"$EXPECTED\", got \"$got\""
fi

# Reshuffle the same lines into a new order and rerun: a hard-coded
# answer or line-position-dependent logic breaks here even though the
# first check above already passed.
shuf "$ORIG" > "$LOG"
got2=$(bash "$SCRIPT" 2>/dev/null | tail -1 | tr -s ' ')
if [ "$got2" != "$EXPECTED" ]; then
  fail "the script's answer changed after the log was reshuffled — it must not depend on line order"
fi

echo "PASS: pipeline correctly identifies $EXPECTED regardless of input order"
