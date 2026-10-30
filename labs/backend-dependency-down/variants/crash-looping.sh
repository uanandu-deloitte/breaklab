#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/frontend deployment/cache deployment/backend --timeout=60s >/dev/null
# Oops: a migration script shipped as the container's startup command and exits immediately
kubectl patch deployment backend --type=strategic -p \
  '{"spec":{"template":{"spec":{"containers":[{"name":"backend","command":["sh","-c","exit 1"]}]}}}}' >/dev/null
kubectl scale deployment backend --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=backend --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment backend --replicas=1 >/dev/null
sleep 5
