#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/orders --timeout=60s >/dev/null
# Oops: the container's port changed in a later edit and the Service was never updated to match
kubectl patch service orders -p '{"spec":{"ports":[{"port":80,"targetPort":9999}]}}' >/dev/null
sleep 2
