#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/worker --timeout=60s >/dev/null
# Oops: a "secret rotation" removed the token and never put a new one back
kubectl set env deployment/worker REQUIRED_TOKEN- >/dev/null
kubectl scale deployment worker --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=worker --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment worker --replicas=1 >/dev/null
sleep 5
