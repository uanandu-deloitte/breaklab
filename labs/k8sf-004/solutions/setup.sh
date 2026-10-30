#!/usr/bin/env bash
set -euo pipefail
pod=$(kubectl get pods -l app=flaky-worker -o jsonpath='{.items[0].metadata.name}')
restarts=$(kubectl get pod "$pod" -o jsonpath='{.status.containerStatuses[0].restartCount}')
reason=$(kubectl get pod "$pod" -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}')
cat <<EOF > /workspace/answers.env
RESTART_COUNT=${restarts}
PREVIOUS_EXIT_REASON=${reason}
EOF
