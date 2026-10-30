#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

CAPTURE=/workspace/capture.pcap
[ -s "$CAPTURE" ] || fail "/workspace/capture.pcap is missing or empty"

read_out=$(tcpdump -r "$CAPTURE" -nn 2>/dev/null) || fail "tcpdump could not read /workspace/capture.pcap"
[ -n "$read_out" ] || fail "capture.pcap contains no packets"

# A request packet has port 8005 as its destination; a response packet
# has port 8005 as its source — so a real two-way exchange needs both
# shapes present, not just any traffic mentioning the port.
echo "$read_out" | grep -Eq '> .*\.8005:' || fail "no packet with destination port 8005 found in the capture — was the request ever sent?"
echo "$read_out" | grep -Eq '\.8005 >' || fail "no packet with source port 8005 found in the capture — was the response ever captured?"

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ "${DEST_PORT:-}" = "8005" ]            || fail "DEST_PORT should be 8005"
[ "${CAPTURE_HAS_REQUEST:-}" = "yes" ]   || fail "CAPTURE_HAS_REQUEST should be yes"
[ "${CAPTURE_HAS_RESPONSE:-}" = "yes" ]  || fail "CAPTURE_HAS_RESPONSE should be yes"

echo "PASS: capture.pcap shows a full request/response exchange on port 8005"
