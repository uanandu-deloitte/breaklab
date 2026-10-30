#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/frontend deployment/cache deployment/backend --timeout=60s >/dev/null
# Oops: a capacity-reduction change scaled the wrong deployment to zero
kubectl scale deployment backend --replicas=0 >/dev/null
sleep 5
