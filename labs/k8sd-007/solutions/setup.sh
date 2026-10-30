#!/usr/bin/env bash
set -euo pipefail

kubectl patch application web -n argocd --type merge \
  -p '{"spec":{"source":{"path":"manifests/web"}}}'
