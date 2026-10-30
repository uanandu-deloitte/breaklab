#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/api --timeout=60s >/dev/null
