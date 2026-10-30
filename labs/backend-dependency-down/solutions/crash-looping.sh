#!/usr/bin/env bash
set -euo pipefail
kubectl rollout undo deployment/backend
kubectl rollout status deployment/backend --timeout=30s
