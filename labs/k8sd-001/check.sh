#!/usr/bin/env bash
set -uo pipefail

pod=$(kubectl get pod -l app=configmap-demo --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1:].metadata.name}' 2>/dev/null)
if [ -z "$pod" ]; then
  echo "FAIL: no configmap-demo pod found"
  exit 1
fi

log_level=$(kubectl get configmap app-config -o jsonpath='{.data.LOG_LEVEL}' 2>/dev/null)
app_conf=$(kubectl get configmap app-config -o jsonpath="{.data['app\.conf']}" 2>/dev/null)

env_val=$(kubectl exec "$pod" -- printenv APP_LOG_LEVEL 2>/dev/null)
file_val=$(kubectl exec "$pod" -- cat /etc/app/app.conf 2>/dev/null)

ok=true
if [ -z "$log_level" ] || [ "$env_val" != "$log_level" ]; then
  echo "FAIL: APP_LOG_LEVEL env var is '${env_val:-<unset>}', ConfigMap's LOG_LEVEL key is '${log_level:-<unset>}'"
  ok=false
fi
if [ -z "$app_conf" ] || [ "$file_val" != "$app_conf" ]; then
  echo "FAIL: /etc/app/app.conf content does not match the ConfigMap's app.conf key"
  ok=false
fi

if [ "$ok" = true ]; then
  echo "PASS: both ConfigMap keys are correctly wired into configmap-demo"
else
  exit 1
fi
