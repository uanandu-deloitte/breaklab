#!/usr/bin/env bash
set -euo pipefail
kubectl rollout undo deployment/api
kubectl rollout status deployment/api --timeout=30s
