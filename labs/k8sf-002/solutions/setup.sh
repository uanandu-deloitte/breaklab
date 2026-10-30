#!/usr/bin/env bash
set -euo pipefail
kubectl set image deployment/web web=nginx:1.27-alpine
kubectl rollout status deployment/web --timeout=120s
