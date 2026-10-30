#!/usr/bin/env bash
set -euo pipefail

kubectl apply -f - <<'YAML'
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
  namespace: limits-demo
spec:
  limits:
    - type: Container
      defaultRequest:
        cpu: 100m
        memory: 64Mi
      default:
        cpu: 200m
        memory: 128Mi
YAML

kubectl apply -f - <<'YAML'
apiVersion: v1
kind: Pod
metadata:
  name: pod-after
  namespace: limits-demo
spec:
  restartPolicy: Never
  containers:
    - name: sleeper
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
YAML
kubectl wait --for=condition=Ready pod/pod-after -n limits-demo --timeout=60s
