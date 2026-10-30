#!/usr/bin/env bash
set -euo pipefail

kubectl apply -f - <<'YAML'
apiVersion: v1
kind: Pod
metadata:
  name: pod-b
  namespace: quota-demo
  labels:
    app: quota-demo
spec:
  restartPolicy: Never
  containers:
    - name: sleeper
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
YAML
kubectl wait --for=condition=Ready pod/pod-b -n quota-demo --timeout=60s

set +e
err=$(kubectl apply -f - 2>&1 <<'YAML'
apiVersion: v1
kind: Pod
metadata:
  name: pod-c
  namespace: quota-demo
  labels:
    app: quota-demo
spec:
  restartPolicy: Never
  containers:
    - name: sleeper
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
YAML
)
set -e

single_line=$(printf '%s' "$err" | tr '\n' ' ' | sed 's/  */ /g')
mkdir -p /workspace
printf 'REJECTION_REASON=%q\n' "$single_line" >> /workspace/answers.env
