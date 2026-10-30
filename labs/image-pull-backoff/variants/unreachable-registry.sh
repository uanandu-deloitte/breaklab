#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/catalog --timeout=60s >/dev/null
# Oops: someone repointed the release at an internal registry that was
# never actually stood up (or isn't reachable from this cluster)
kubectl set image deployment/catalog catalog=registry.internal.test:5000/catalog:v2 >/dev/null
kubectl rollout status deployment/catalog --timeout=20s >/dev/null 2>&1 || true
# Force out the old, still-Ready pod (see bad-tag.sh for why a rolling
# update alone leaves it running).
kubectl scale deployment catalog --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=catalog --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment catalog --replicas=1 >/dev/null
sleep 5
