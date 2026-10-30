#!/usr/bin/env bash
set -euo pipefail
kubectl rollout undo deployment/worker
kubectl rollout status deployment/worker --timeout=30s
