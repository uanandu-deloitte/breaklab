#!/usr/bin/env bash
set -uo pipefail
ready=$(kubectl get pod oom-demo -o jsonpath='{.status.containerStatuses[0].ready}' 2>/dev/null)
sleep 8
ready2=$(kubectl get pod oom-demo -o jsonpath='{.status.containerStatuses[0].ready}' 2>/dev/null)
if [ "$ready" = "true" ] && [ "$ready2" = "true" ]; then
  echo "PASS: oom-demo has stayed up"
else
  reason=$(kubectl get pod oom-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}' 2>/dev/null)
  echo "FAIL: oom-demo is not staying up (last termination: ${reason:-unknown})"
  exit 1
fi
