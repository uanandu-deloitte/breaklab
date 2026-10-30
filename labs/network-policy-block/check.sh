#!/usr/bin/env bash
set -uo pipefail
pod=$(kubectl get pod -l app=frontend -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
for i in $(seq 1 10); do
  if kubectl exec "$pod" -- wget -q -T 2 -O- http://cache:6379/ >/dev/null 2>&1; then
    echo "PASS: frontend reports the cache as reachable"
    exit 0
  fi
  sleep 1
done
echo "FAIL: frontend's most recent cache check was 'cache: unreachable'"
exit 1
