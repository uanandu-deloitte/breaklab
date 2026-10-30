#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

live_phase=$(kubectl get pvc reports-pvc -o jsonpath='{.status.phase}' 2>/dev/null)
live_requested=$(kubectl get pvc reports-pvc -o jsonpath='{.spec.storageClassName}' 2>/dev/null)
live_available=$(kubectl get storageclass -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n' | sort | tr '\n' ' ' | sed 's/ *$//')
live_pod_phase=$(kubectl get pod reports-writer -o jsonpath='{.status.phase}' 2>/dev/null)

[ -n "${PVC_PHASE:-}" ] || fail "PVC_PHASE is empty"
[ "$PVC_PHASE" = "$live_phase" ] || fail "PVC_PHASE ($PVC_PHASE) doesn't match reports-pvc's live phase ($live_phase)"

[ -n "${REQUESTED_STORAGECLASS:-}" ] || fail "REQUESTED_STORAGECLASS is empty"
[ "$REQUESTED_STORAGECLASS" = "$live_requested" ] || fail "REQUESTED_STORAGECLASS ($REQUESTED_STORAGECLASS) doesn't match reports-pvc's live spec.storageClassName ($live_requested)"

[ -n "${AVAILABLE_STORAGECLASSES:-}" ] || fail "AVAILABLE_STORAGECLASSES is empty"
normalized_answer=$(echo "$AVAILABLE_STORAGECLASSES" | tr ' ' '\n' | sort | tr '\n' ' ' | sed 's/ *$//')
[ "$normalized_answer" = "$live_available" ] || fail "AVAILABLE_STORAGECLASSES ('$AVAILABLE_STORAGECLASSES') doesn't match what kubectl get storageclass actually lists ('$live_available')"

[ -n "${CLASS_EXISTS:-}" ] || fail "CLASS_EXISTS is empty"
exists="no"
for sc in $live_available; do
  [ "$sc" = "$live_requested" ] && exists="yes"
done
echo "$CLASS_EXISTS" | grep -qi "^${exists}$" || fail "CLASS_EXISTS ($CLASS_EXISTS) should be '$exists' — requested class '$live_requested' against available classes '$live_available'"

[ -n "${POD_PHASE:-}" ] || fail "POD_PHASE is empty"
[ "$POD_PHASE" = "$live_pod_phase" ] || fail "POD_PHASE ($POD_PHASE) doesn't match reports-writer's live phase ($live_pod_phase)"

echo "PASS: correctly diagnosed reports-pvc as $live_phase — requested '$live_requested', available '$live_available'"
