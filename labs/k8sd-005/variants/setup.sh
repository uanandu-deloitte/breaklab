#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
for i in $(seq 1 30); do
  reason=$(kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}' 2>/dev/null)
  [ "$reason" = "OOMKilled" ] && break
  sleep 2
done
