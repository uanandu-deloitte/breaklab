#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/frontend deployment/cache deployment/backend --timeout=60s >/dev/null
# Oops: a config change shipped with a typo'd startup command
kubectl patch deployment cache --type=strategic -p \
  '{"spec":{"template":{"spec":{"containers":[{"name":"cache","command":["sh","-c","exit 1"]}]}}}}' >/dev/null
# Force out the old, still-Ready pod (a rolling update alone won't
# remove it until a replacement is Ready, which never happens once the
# new template crashes on start).
kubectl scale deployment cache --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=cache --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment cache --replicas=1 >/dev/null
sleep 5
