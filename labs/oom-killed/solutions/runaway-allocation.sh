#!/usr/bin/env bash
set -euo pipefail
kubectl delete pod oom-demo --wait=true
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: Pod
metadata: {name: oom-demo}
spec:
  restartPolicy: Always
  containers:
  - name: oom-demo
    image: python:3.12-alpine
    command: ["python","-c","a=[]; [a.append(bytearray(1024*1024)) for _ in range(5)]; import time; time.sleep(3600)"]
    resources:
      requests: {memory: 8Mi}
      limits: {memory: 24Mi}
YAML
kubectl wait --for=condition=Ready pod/oom-demo --timeout=30s
