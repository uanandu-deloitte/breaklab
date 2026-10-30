#!/usr/bin/env bash
set -uo pipefail

RESULT=/tmp/last-request-result.txt
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$RESULT" ] || fail "no request has reached the fixture server yet — send one with curl first"

outcome=$(cat "$RESULT")
[ "$outcome" = "success" ] || fail "the last request the server saw did not succeed: $outcome"

echo "PASS: the fixture server recorded a successful request"
