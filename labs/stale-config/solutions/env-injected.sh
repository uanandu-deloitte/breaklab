#!/usr/bin/env bash
set -euo pipefail
# Reference fix: env vars are injected once at container start, so the only
# way to pick up the new ConfigMap value is a fresh pod
kubectl rollout restart deployment/configwatcher
kubectl rollout status deployment/configwatcher --timeout=60s
