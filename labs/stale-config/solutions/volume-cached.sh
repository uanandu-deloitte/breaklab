#!/usr/bin/env bash
set -euo pipefail
# Reference fix: the mounted file is already current, but the app cached it
# at startup, so the only way to make it re-read is a fresh pod
kubectl rollout restart deployment/configwatcher
kubectl rollout status deployment/configwatcher --timeout=60s
