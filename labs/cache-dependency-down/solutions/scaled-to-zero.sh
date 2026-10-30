#!/usr/bin/env bash
set -euo pipefail
kubectl scale deployment cache --replicas=1
kubectl rollout status deployment/cache --timeout=30s
