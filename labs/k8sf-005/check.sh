#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

live_ready_column=$(kubectl get pods -l app=catalog-api --no-headers 2>/dev/null | awk '{print $2}')
[ -n "$live_ready_column" ] || fail "catalog-api Pod not found"

[ -n "${READY_COLUMN:-}" ] || fail "READY_COLUMN is empty"
[ "$READY_COLUMN" = "$live_ready_column" ] || fail "READY_COLUMN ($READY_COLUMN) doesn't match the Pod's live READY column ($live_ready_column)"

live_addrs=$(kubectl get endpoints catalog-api -o jsonpath='{.subsets[0].addresses[*].ip}' 2>/dev/null)
if [ -z "$live_addrs" ]; then
  live_count=0
else
  live_count=$(echo "$live_addrs" | wc -w | tr -d ' ')
fi

[ -n "${ENDPOINT_COUNT:-}" ] || fail "ENDPOINT_COUNT is empty"
[ "$ENDPOINT_COUNT" = "$live_count" ] || fail "ENDPOINT_COUNT ($ENDPOINT_COUNT) doesn't match the Service's live endpoint count ($live_count)"
[ "$live_count" = "0" ] || fail "expected the catalog-api Service to have zero live endpoints, found $live_count"

live_probe_path=$(kubectl get deployment catalog-api -o jsonpath='{.spec.template.spec.containers[0].readinessProbe.httpGet.path}' 2>/dev/null)
[ -n "$live_probe_path" ] || fail "could not read the live readinessProbe path"

[ -n "${PROBE_PATH:-}" ] || fail "PROBE_PATH is empty"
[ "$PROBE_PATH" = "$live_probe_path" ] || fail "PROBE_PATH ($PROBE_PATH) doesn't match the Deployment's live readinessProbe path ($live_probe_path)"

echo "PASS: recorded readiness state, endpoint count, and probe path all match live cluster state"
