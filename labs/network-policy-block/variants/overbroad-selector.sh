#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/frontend deployment/cache deployment/backend --timeout=60s >/dev/null
sleep 5
# Oops: this was meant to lock down just the cache Pods, but the
# podSelector was left empty, so it matches every Pod in the namespace
cat > /tmp/lockdown.yaml <<'YAML'
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: {name: lockdown-cache}
spec:
  podSelector: {}
  policyTypes: [Ingress]
  ingress: []
YAML
kubectl apply -f /tmp/lockdown.yaml >/dev/null
sleep 5
