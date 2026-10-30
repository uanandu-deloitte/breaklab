#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/storefront --timeout=60s >/dev/null
# Oops: a probe "fix" pointed readiness at a file that's never created
kubectl patch deployment storefront --type=json \
  -p='[{"op":"replace","path":"/spec/template/spec/containers/0/readinessProbe/exec/command","value":["cat","/tmp/never-exists"]}]' >/dev/null
# A rolling update alone won't remove the old, still-Ready pod: with
# replicas:1 the Deployment controller won't scale down the old
# ReplicaSet until a replacement is Ready, which never happens once the
# probe is broken. Cycle replicas to 0 and back so only the new (broken)
# template is ever running.
kubectl scale deployment storefront --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=storefront --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment storefront --replicas=1 >/dev/null
sleep 5
