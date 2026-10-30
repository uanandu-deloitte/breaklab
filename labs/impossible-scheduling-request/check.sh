#!/usr/bin/env bash
set -uo pipefail
phase=$(kubectl get pod worker-task -o jsonpath='{.status.phase}' 2>/dev/null)
if [ "$phase" = "Running" ]; then
  echo "PASS: worker-task is Running"
else
  echo "FAIL: worker-task is ${phase:-missing}, not Running"
  exit 1
fi
