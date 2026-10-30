#!/usr/bin/env bash
set -euo pipefail
# Reference fix: roll the Deployment back to its working readiness probe
kubectl rollout undo deployment/checkout
kubectl rollout status deployment/checkout --timeout=30s
# Reference fix: restore the correct comparison direction
sed -i 's/^op: .*/op: "<"/' /etc/alerts/rules.yaml
