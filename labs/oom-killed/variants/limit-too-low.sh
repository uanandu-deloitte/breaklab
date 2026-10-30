#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
# Oops: this limit was copied from a much smaller sibling service
kubectl apply -f - <<'YAML' >/dev/null
apiVersion: v1
kind: Pod
metadata: {name: oom-demo}
spec:
  restartPolicy: Always
  containers:
  - name: oom-demo
    image: python:3.12-alpine
    command: ["python","-c","a=[]; [a.append(bytearray(1024*1024)) for _ in range(10)]; import time; time.sleep(3600)"]
    resources:
      requests: {memory: 4Mi}
      limits: {memory: 6Mi}
YAML
sleep 10
