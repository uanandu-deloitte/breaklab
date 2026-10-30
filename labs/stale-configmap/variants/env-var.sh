#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/config-reader --timeout=60s >/dev/null
# Oops: the ConfigMap was updated to v2, but nothing restarted the Deployment
kubectl create configmap app-config \
  --from-literal=APP_MESSAGE=v2 --from-literal=message=v1 \
  -o yaml --dry-run=client | kubectl apply -f - >/dev/null
sleep 2
