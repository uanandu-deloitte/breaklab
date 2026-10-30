#!/usr/bin/env bash
set -euo pipefail
sed -i 's/validationFailureAction: Audit/validationFailureAction: Enforce/' /manifests/policy.yaml
kubectl apply -f /manifests/policy.yaml
sleep 3
kubectl delete pod bad-pod --wait=true
