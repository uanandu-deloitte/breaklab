#!/usr/bin/env bash
set -euo pipefail
kubectl rollout undo deployment/cache
kubectl rollout status deployment/cache --timeout=30s
