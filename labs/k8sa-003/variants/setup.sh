#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
# Oops: this node got tainted for a dedicated workload pool, but the
# existing manifest was never updated to tolerate it.
node=$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')
kubectl taint nodes "$node" dedicated=special:NoSchedule --overwrite >/dev/null
kubectl apply -f /manifests/ >/dev/null
sleep 3
