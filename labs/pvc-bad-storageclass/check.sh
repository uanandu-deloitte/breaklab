#!/usr/bin/env bash
set -uo pipefail
phase=$(kubectl get pvc data-pending -o jsonpath='{.status.phase}' 2>/dev/null)
if [ "$phase" = "Bound" ]; then
  echo "PASS: data-pending is Bound"
else
  echo "FAIL: data-pending is ${phase:-missing}, not Bound"
  exit 1
fi
