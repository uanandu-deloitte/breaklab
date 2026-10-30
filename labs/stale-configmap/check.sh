#!/usr/bin/env bash
set -uo pipefail
pod=$(kubectl get pod -l app=config-reader -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
if [ -z "$pod" ]; then
  echo "FAIL: config-reader Pod not found"
  exit 1
fi
env_val=$(kubectl exec "$pod" -- printenv APP_MESSAGE 2>/dev/null || true)
file_val=$(kubectl exec "$pod" -- cat /etc/config/message 2>/dev/null || true)
if [ "$env_val" = "v2" ] || [ "$file_val" = "v2" ]; then
  echo "PASS: config-reader sees the current value (v2)"
else
  echo "FAIL: config-reader still sees env='${env_val:-unset}' file='${file_val:-unset}', not v2"
  exit 1
fi
