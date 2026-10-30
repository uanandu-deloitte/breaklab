#!/usr/bin/env bash
set -uo pipefail

ANSWER=/workspace/answer.txt
EXPECTED=/var/log/app/archive/2024-01/host-03/trace.log
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWER" ] || fail "/workspace/answer.txt not found"

got=$(sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' "$ANSWER")

[ "$got" = "$EXPECTED" ] || fail "expected \"$EXPECTED\", got \"$got\""

[ -f "$EXPECTED" ] || fail "the real log file $EXPECTED is missing — fixture was tampered with"
grep -q "transaction-id=TX-8842" "$EXPECTED" || fail "the real log file no longer contains the transaction ID"

echo "PASS: answer.txt names the log file that actually contains TX-8842"
