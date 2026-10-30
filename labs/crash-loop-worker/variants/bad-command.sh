#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/worker --timeout=60s >/dev/null
# Oops: a "fix" to the entrypoint shipped with a typo and the process exits immediately
kubectl patch deployment worker --type=strategic -p \
  '{"spec":{"template":{"spec":{"containers":[{"name":"worker","command":["sh","-c","echo intentional-failure >&2; exit 1"]}]}}}}' >/dev/null
kubectl scale deployment worker --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=worker --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment worker --replicas=1 >/dev/null
sleep 5
