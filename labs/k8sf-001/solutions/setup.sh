#!/usr/bin/env bash
set -euo pipefail
kubectl scale deployment/hello --replicas=4
kubectl wait --for=jsonpath='{.status.readyReplicas}'=4 deployment/hello --timeout=60s
image=$(kubectl get deployment hello -o jsonpath='{.spec.template.spec.containers[0].image}')
initial=$(awk -F': *' '/replicas:/{print $2; exit}' /manifests/app.yaml)
cat <<EOF > /workspace/answers.env
INITIAL_REPLICA_COUNT=${initial}
CONTAINER_IMAGE=${image}
FINAL_REPLICA_COUNT=4
EOF
