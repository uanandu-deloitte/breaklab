#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

ns=quota-demo
quota_name=pod-quota

used=$(kubectl get resourcequota "$quota_name" -n "$ns" -o jsonpath='{.status.used.pods}' 2>/dev/null)
hard=$(kubectl get resourcequota "$quota_name" -n "$ns" -o jsonpath='{.status.hard.pods}' 2>/dev/null)
[ -n "$hard" ] || fail "resourcequota $quota_name not found in namespace $ns"
[ "$hard" = "2" ] || fail "resourcequota hard pod limit is ${hard:-unset}, expected 2"
[ "$used" = "$hard" ] || fail "resourcequota shows ${used:-0}/${hard} pods used — it hasn't reached its limit yet"

running=$(kubectl get pods -n "$ns" --field-selector=status.phase=Running --no-headers 2>/dev/null | wc -l | tr -d ' ')
[ "$running" = "2" ] || fail "expected exactly 2 Running pods in $ns, found $running"

answers=/workspace/answers.env
[ -f "$answers" ] || fail "/workspace/answers.env not found — record the quota rejection reason there"

# shellcheck disable=SC1090
source "$answers"
[ -n "${REJECTION_REASON:-}" ] || fail "REJECTION_REASON is empty in answers.env"
echo "$REJECTION_REASON" | grep -qiE "exceeded quota|forbidden" || \
  fail "REJECTION_REASON doesn't look like a real quota rejection: $REJECTION_REASON"

echo "PASS: quota-demo is at its 2/2 pod quota and the rejection reason was recorded"
