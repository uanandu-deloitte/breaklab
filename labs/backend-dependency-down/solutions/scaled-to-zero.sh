#!/usr/bin/env bash
set -euo pipefail
kubectl scale deployment backend --replicas=1
kubectl rollout status deployment/backend --timeout=30s
