#!/usr/bin/env bash
set -uo pipefail
count=$(kubectl get endpoints api -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null | wc -w)
if [ "$count" -ge 1 ]; then
  echo "PASS: api Service has $count healthy backend(s)"
else
  echo "FAIL: api Service has no healthy backends"
  exit 1
fi
