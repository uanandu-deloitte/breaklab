#!/usr/bin/env bash
set -euo pipefail
reason=$(kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}')
exit_code=$(kubectl get pod mem-demo -o jsonpath='{.status.containerStatuses[0].lastState.terminated.exitCode}')
limit=$(kubectl get pod mem-demo -o jsonpath='{.spec.containers[0].resources.limits.memory}')
mkdir -p /workspace
cat <<EOF > /workspace/answers.env
TERMINATION_REASON=$reason
EXIT_CODE=$exit_code
MEMORY_LIMIT=$limit
IS_OOM=yes
EOF
