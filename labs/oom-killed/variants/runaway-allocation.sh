#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
# Oops: a new code path appends data to an in-memory list and never bounds it
kubectl apply -f - <<'YAML' >/dev/null
apiVersion: v1
kind: Pod
metadata: {name: oom-demo}
spec:
  restartPolicy: Always
  containers:
  - name: oom-demo
    image: python:3.12-alpine
    command: ["python","-c","a=[]; [a.append(bytearray(1024*1024)) for _ in range(256)]"]
    resources:
      requests: {memory: 8Mi}
      limits: {memory: 24Mi}
YAML
sleep 10
