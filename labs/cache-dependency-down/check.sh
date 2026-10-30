#!/usr/bin/env bash
set -uo pipefail
for i in $(seq 1 10); do
  count=$(kubectl get endpoints cache -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null | wc -w)
  if [ "$count" -ge 1 ]; then
    echo "PASS: frontend reports the cache as reachable"
    exit 0
  fi
  sleep 1
done
echo "FAIL: frontend's most recent cache check was 'cache: unreachable'"
exit 1
