#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null

# catalog-api's readiness probe checks a path the container never
# serves, so the Deployment will never reach condition=Available — wait
# for the Pod to exist and be Running instead of waiting on readiness.
for i in $(seq 1 60); do
  phase=$(kubectl get pods -l app=catalog-api -o jsonpath='{.items[0].status.phase}' 2>/dev/null || true)
  [ "$phase" = "Running" ] && exit 0
  sleep 1
done

echo "catalog-api Pod did not reach Running in time" >&2
exit 1
