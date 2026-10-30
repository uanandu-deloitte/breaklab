#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/frontend deployment/cache deployment/backend --timeout=60s >/dev/null
# Let cross-service traffic settle before breaking it, so the "healthy
# baseline" isn't confused with this lab's actual fault.
sleep 5
# Oops: a security pass locked down cache with no matching allow rule
cat > /tmp/deny-cache.yaml <<'YAML'
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: {name: lockdown-cache}
spec:
  podSelector: {matchLabels: {app: cache}}
  policyTypes: [Ingress]
  ingress: []
YAML
kubectl apply -f /tmp/deny-cache.yaml >/dev/null
sleep 5
