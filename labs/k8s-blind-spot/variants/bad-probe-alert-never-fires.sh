#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/ >/dev/null
kubectl wait --for=condition=Available deployment/checkout --timeout=60s >/dev/null

# Oops (k8s side): the readiness probe was repointed at a port nothing
# listens on, so the real Pod never becomes Ready
kubectl patch deployment checkout --type=json \
  -p='[{"op":"replace","path":"/spec/template/spec/containers/0/readinessProbe/httpGet/port","value":9999}]' >/dev/null
kubectl scale deployment checkout --replicas=0 >/dev/null
kubectl wait --for=delete pod -l app=checkout --timeout=30s >/dev/null 2>&1 || true
kubectl scale deployment checkout --replicas=1 >/dev/null

# Oops (observability side): a "tune the threshold" pass flipped the
# comparison direction, so the rule can never fire no matter the value
sed -i 's/^op: .*/op: ">"/' /etc/alerts/rules.yaml
sleep 5
