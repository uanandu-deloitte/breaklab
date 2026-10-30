#!/usr/bin/env bash
set -uo pipefail
count=$(kubectl get endpoints storefront -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null | wc -w)
if [ "$count" -ge 1 ]; then
  echo "PASS: storefront Service has $count healthy backend(s)"
else
  echo "FAIL: storefront Service has no healthy backends"
  exit 1
fi
