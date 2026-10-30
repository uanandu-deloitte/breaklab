#!/usr/bin/env bash
set -uo pipefail
pod=$(kubectl get pods -l app=configwatcher --no-headers 2>/dev/null | grep -v Terminating | awk 'NR==1{print $1}')
if [ -z "$pod" ]; then
  echo "FAIL: no configwatcher pod found"
  exit 1
fi
status=$(kubectl exec "$pod" -- cat /tmp/status 2>/dev/null | tr -d '[:space:]')
if [ "$status" = "normal" ]; then
  echo "PASS: configwatcher reports mode=normal"
else
  echo "FAIL: configwatcher still reports mode=${status:-unknown}"
  exit 1
fi
