#!/usr/bin/env bash
set -uo pipefail
ready=$(kubectl get deployment worker -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
if [ "$ready" = "1" ]; then
  echo "PASS: worker Deployment has a ready replica"
else
  echo "FAIL: worker Deployment has no ready replica (readyReplicas=${ready:-0})"
  exit 1
fi
