#!/usr/bin/env bash
set -euo pipefail
for i in $(seq 1 90); do
  kubectl get nodes --no-headers 2>/dev/null | grep -q . && break
  sleep 1
done
kubectl wait --for=condition=Ready node --all --timeout=60s >/dev/null
for i in $(seq 1 60); do
  kubectl get serviceaccount default >/dev/null 2>&1 && break
  sleep 1
done
