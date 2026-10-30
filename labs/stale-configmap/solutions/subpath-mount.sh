#!/usr/bin/env bash
set -euo pipefail
kubectl rollout restart deployment/config-reader
kubectl rollout status deployment/config-reader --timeout=30s
