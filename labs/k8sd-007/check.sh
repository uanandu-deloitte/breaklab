#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

sync_status=""
health_status=""
for i in $(seq 1 30); do
  sync_status=$(kubectl get application web -n argocd -o jsonpath='{.status.sync.status}' 2>/dev/null)
  health_status=$(kubectl get application web -n argocd -o jsonpath='{.status.health.status}' 2>/dev/null)
  [ "$sync_status" = "Synced" ] && [ "$health_status" = "Healthy" ] && break
  sleep 2
done

[ "$sync_status" = "Synced" ] || fail "Application web is not Synced (status: ${sync_status:-unknown})"
[ "$health_status" = "Healthy" ] || fail "Application web is not Healthy (status: ${health_status:-unknown})"

ready=$(kubectl get deployment web -n default -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
[ "${ready:-0}" -ge 1 ] || fail "Deployment web in namespace default has no ready replicas -- ArgoCD reports Synced but the real workload isn't actually running"

echo "PASS: Application web is Synced and Healthy, and the real web Deployment is running"
