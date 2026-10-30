#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/storefront --timeout=60s >/dev/null
# Oops: a "cleanup" pass retyped the service's selector
kubectl patch service storefront --type merge -p '{"spec":{"selector":{"app":"store-front"}}}' >/dev/null
