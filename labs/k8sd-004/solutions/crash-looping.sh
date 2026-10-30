#!/usr/bin/env bash
set -euo pipefail
dep="backend"
endpoints=$(kubectl get endpoints "$dep" -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null | wc -w | tr -d ' ')
desired=$(kubectl get deployment "$dep" -o jsonpath='{.spec.replicas}' 2>/dev/null)
ready=$(kubectl get deployment "$dep" -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
mkdir -p /workspace
cat <<EOF > /workspace/answers.env
DEPENDENCY_NAME=$dep
READY_ENDPOINTS=${endpoints:-0}
DESIRED_REPLICAS=${desired:-0}
READY_REPLICAS=${ready:-0}
FAILURE_MODE=crash-looping
EOF
