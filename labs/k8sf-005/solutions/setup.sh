#!/usr/bin/env bash
set -euo pipefail
ready_column=$(kubectl get pods -l app=catalog-api --no-headers | awk '{print $2}')
addrs=$(kubectl get endpoints catalog-api -o jsonpath='{.subsets[0].addresses[*].ip}' 2>/dev/null || true)
count=0
[ -n "$addrs" ] && count=$(echo "$addrs" | wc -w | tr -d ' ')
probe_path=$(kubectl get deployment catalog-api -o jsonpath='{.spec.template.spec.containers[0].readinessProbe.httpGet.path}')
cat <<EOF > /workspace/answers.env
READY_COLUMN=${ready_column}
ENDPOINT_COUNT=${count}
PROBE_PATH=${probe_path}
EOF
