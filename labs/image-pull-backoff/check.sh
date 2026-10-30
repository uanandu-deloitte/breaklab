#!/usr/bin/env bash
set -uo pipefail
ready=$(kubectl get deployment catalog -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
if [ "$ready" = "1" ]; then
  echo "PASS: catalog Deployment has a ready replica"
else
  echo "FAIL: catalog Deployment has no ready replica (readyReplicas=${ready:-0})"
  exit 1
fi
