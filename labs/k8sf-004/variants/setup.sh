#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null

# The flaky-worker Pod crashes twice (via a counter file in an emptyDir
# volume, so the count survives container restarts within the same Pod)
# before it stabilizes on its third attempt. Wait for restartCount to
# reach 2 and the container to actually be running, rather than racing
# the learner against crash backoff delays.
for i in $(seq 1 120); do
  pod=$(kubectl get pods -l app=flaky-worker -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)
  if [ -n "$pod" ]; then
    restarts=$(kubectl get pod "$pod" -o jsonpath='{.status.containerStatuses[0].restartCount}' 2>/dev/null || echo "")
    phase=$(kubectl get pod "$pod" -o jsonpath='{.status.phase}' 2>/dev/null || echo "")
    if [ "${restarts:-0}" -ge 2 ] 2>/dev/null && [ "$phase" = "Running" ]; then
      exit 0
    fi
  fi
  sleep 2
done

echo "flaky-worker did not stabilize in time" >&2
exit 1
