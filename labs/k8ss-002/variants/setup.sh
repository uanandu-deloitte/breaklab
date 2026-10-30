#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Ready pod/scratch-writer --timeout=60s >/dev/null
sleep 2
