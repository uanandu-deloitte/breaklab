#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point the readiness probe back at the file the app actually creates
kubectl patch deployment storefront --type=json \
  -p='[{"op":"replace","path":"/spec/template/spec/containers/0/readinessProbe/exec/command","value":["cat","/tmp/ready"]}]'
kubectl rollout status deployment/storefront --timeout=60s
