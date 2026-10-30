#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/frontend deployment/cache deployment/backend --timeout=60s >/dev/null
# Oops: someone scaled cache down during a "cleanup" and never scaled it back
kubectl scale deployment cache --replicas=0 >/dev/null
sleep 5
