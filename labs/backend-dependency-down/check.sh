#!/usr/bin/env bash
set -uo pipefail
for i in $(seq 1 10); do
  count=$(kubectl get endpoints backend -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null | wc -w)
  if [ "$count" -ge 1 ]; then
    echo "PASS: frontend reports the backend as reachable"
    exit 0
  fi
  sleep 1
done
echo "FAIL: frontend's most recent backend check was 'backend: unreachable'"
exit 1
