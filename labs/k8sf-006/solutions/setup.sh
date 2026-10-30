#!/usr/bin/env bash
set -euo pipefail
selector=$(kubectl get svc reports -o jsonpath='{.spec.selector.app}')
pod_label=$(kubectl get pods -l tier=backend -o jsonpath='{.items[0].metadata.labels.app}')
addrs=$(kubectl get endpoints reports -o jsonpath='{.subsets[0].addresses[*].ip}' 2>/dev/null || true)
count=0
[ -n "$addrs" ] && count=$(echo "$addrs" | wc -w | tr -d ' ')
cat <<EOF > /workspace/answers.env
SERVICE_SELECTOR_VALUE=${selector}
POD_LABEL_VALUE=${pod_label}
ENDPOINT_COUNT=${count}
EOF
