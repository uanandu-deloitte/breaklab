#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/api --timeout=60s >/dev/null
# Oops: the probe was copied from a different service and never updated to this one's port
kubectl patch deployment api --type=json \
  -p='[{"op":"replace","path":"/spec/template/spec/containers/0/readinessProbe/httpGet/port","value":9999}]' >/dev/null
kubectl scale deployment api --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=api --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment api --replicas=1 >/dev/null
sleep 5
