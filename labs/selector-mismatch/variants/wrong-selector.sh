#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/orders --timeout=60s >/dev/null
# Oops: a "cleanup" pass retyped the Service's selector
kubectl patch service orders -p '{"spec":{"selector":{"app":"order"}}}' >/dev/null
sleep 2
