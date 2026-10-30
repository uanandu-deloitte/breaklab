#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/catalog --timeout=60s >/dev/null
# Oops: the release pointed at a tag that was never actually pushed
kubectl set image deployment/catalog catalog=busybox:no-such-tag >/dev/null
kubectl rollout status deployment/catalog --timeout=20s >/dev/null 2>&1 || true
# A rolling update alone won't remove the old, still-Ready pod: with
# replicas:1 the Deployment controller won't scale down the old
# ReplicaSet until a replacement is Ready, which never happens once the
# image can't be pulled. Cycle replicas to 0 and back so only the new
# (broken) template is ever running.
kubectl scale deployment catalog --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=catalog --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment catalog --replicas=1 >/dev/null
sleep 5
