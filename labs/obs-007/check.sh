#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

command -v log-event-a >/dev/null || fail "log-event-a is missing"
command -v log-event-b >/dev/null || fail "log-event-b is missing"

real_before=$(date +%s)
log-event-a "check" >/dev/null 2>&1
log-event-b "check" >/dev/null 2>&1
real_after=$(date +%s)
real_now=$(( (real_before + real_after) / 2 ))

ts_a=$(tail -n 1 /var/log/host-a.log | awk '{print $2}')
ts_b=$(tail -n 1 /var/log/host-b.log | awk '{print $2}')
case "$ts_a" in ''|*[!0-9]*) fail "could not read a numeric timestamp from /var/log/host-a.log" ;; esac
case "$ts_b" in ''|*[!0-9]*) fail "could not read a numeric timestamp from /var/log/host-b.log" ;; esac

skew_a=$(( ts_a - real_now ))
skew_b=$(( ts_b - real_now ))
abs_skew_a=$skew_a; [ "$abs_skew_a" -lt 0 ] && abs_skew_a=$(( -abs_skew_a ))
abs_skew_b=$skew_b; [ "$abs_skew_b" -lt 0 ] && abs_skew_b=$(( -abs_skew_b ))

if [ "$abs_skew_a" -gt "$abs_skew_b" ]; then
  ground_truth_drifted=host-a
  ground_truth_skew=$abs_skew_a
else
  ground_truth_drifted=host-b
  ground_truth_skew=$abs_skew_b
fi

min_skew=$abs_skew_a
[ "$abs_skew_b" -lt "$min_skew" ] && min_skew=$abs_skew_b
gap=$(( ground_truth_skew - min_skew ))
[ "$gap" -ge 10 ] || fail "both hosts' clocks currently look too close together (host-a skew=${abs_skew_a}s, host-b skew=${abs_skew_b}s) to confidently call one of them drifted — the lab environment may be unhealthy"

[ -n "${DRIFTED_HOST:-}" ] || fail "DRIFTED_HOST is empty"
[ "$DRIFTED_HOST" = "$ground_truth_drifted" ] \
  || fail "DRIFTED_HOST is '$DRIFTED_HOST', but host-a is off by ${abs_skew_a}s and host-b is off by ${abs_skew_b}s from real time — expected '$ground_truth_drifted'"

[ -n "${SKEW_SECONDS:-}" ] || fail "SKEW_SECONDS is empty"
case "$SKEW_SECONDS" in ''|*[!0-9]*) fail "SKEW_SECONDS is not a plain non-negative number: $SKEW_SECONDS" ;; esac
diff_from_truth=$(( SKEW_SECONDS - ground_truth_skew ))
[ "$diff_from_truth" -lt 0 ] && diff_from_truth=$(( -diff_from_truth ))
[ "$diff_from_truth" -le 5 ] \
  || fail "SKEW_SECONDS ($SKEW_SECONDS) is too far from the actual measured skew (~${ground_truth_skew}s) for $ground_truth_drifted"

[ -n "${IS_REAL_ORDERING_ISSUE:-}" ] || fail "IS_REAL_ORDERING_ISSUE is empty"
echo "$IS_REAL_ORDERING_ISSUE" | grep -qi '^no$' \
  || fail "IS_REAL_ORDERING_ISSUE should be 'no' — this is clock drift on $ground_truth_drifted, not a genuine event-ordering problem"

echo "PASS: correctly identified $ground_truth_drifted as drifted by ~${ground_truth_skew}s, not a real ordering issue"
