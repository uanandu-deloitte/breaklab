#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

live_selector=$(kubectl get svc reports -o jsonpath='{.spec.selector.app}' 2>/dev/null)
[ -n "$live_selector" ] || fail "reports Service not found or has no selector.app"

[ -n "${SERVICE_SELECTOR_VALUE:-}" ] || fail "SERVICE_SELECTOR_VALUE is empty"
[ "$SERVICE_SELECTOR_VALUE" = "$live_selector" ] || fail "SERVICE_SELECTOR_VALUE ($SERVICE_SELECTOR_VALUE) doesn't match the Service's live selector.app ($live_selector)"

live_pod_label=$(kubectl get pods -l tier=backend -o jsonpath='{.items[0].metadata.labels.app}' 2>/dev/null)
[ -n "$live_pod_label" ] || fail "could not read the reports Pod's live app label"

[ -n "${POD_LABEL_VALUE:-}" ] || fail "POD_LABEL_VALUE is empty"
[ "$POD_LABEL_VALUE" = "$live_pod_label" ] || fail "POD_LABEL_VALUE ($POD_LABEL_VALUE) doesn't match the Pod's live app label ($live_pod_label)"

[ "$live_selector" != "$live_pod_label" ] || fail "expected the Service's selector and the Pod's label to actually differ in this lab"

live_addrs=$(kubectl get endpoints reports -o jsonpath='{.subsets[0].addresses[*].ip}' 2>/dev/null)
if [ -z "$live_addrs" ]; then
  live_count=0
else
  live_count=$(echo "$live_addrs" | wc -w | tr -d ' ')
fi

[ -n "${ENDPOINT_COUNT:-}" ] || fail "ENDPOINT_COUNT is empty"
[ "$ENDPOINT_COUNT" = "$live_count" ] || fail "ENDPOINT_COUNT ($ENDPOINT_COUNT) doesn't match the Service's live endpoint count ($live_count)"
[ "$live_count" = "0" ] || fail "expected the reports Service to have zero live endpoints, found $live_count"

echo "PASS: recorded selector/label values and endpoint count all match live cluster state, confirming the mismatch"
