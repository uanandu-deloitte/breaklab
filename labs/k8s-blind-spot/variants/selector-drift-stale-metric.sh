#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/checkout --timeout=60s >/dev/null

# Oops (k8s side): a "cleanup" pass retyped the Service's selector
kubectl patch service checkout -p '{"spec":{"selector":{"app":"checkout-svc"}}}' >/dev/null

# Oops (observability side): the periodic refresh wrapper was left
# pointing at a renamed binary from an old rollout, so it has silently
# stopped updating current.json since before the outage. current.json
# still holds whatever was last observed, long ago.
sed -i 's#/usr/local/bin/metrics-agent#/usr/local/bin/metrics-agent-legacy#' /usr/local/bin/refresh-metrics
sleep 2
