#!/usr/bin/env bash
set -euo pipefail
kubectl delete pod worker-task --wait=true
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: Pod
metadata: {name: worker-task}
spec:
  restartPolicy: Never
  containers:
  - name: sleeper
    image: busybox:1.36
    command: ["sh","-c","sleep 3600"]
    resources:
      requests: {cpu: 10m, memory: 16Mi}
YAML
kubectl wait --for=condition=Ready pod/worker-task --timeout=30s
