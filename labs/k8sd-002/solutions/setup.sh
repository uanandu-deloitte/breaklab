#!/usr/bin/env bash
set -euo pipefail
kubectl delete pod qos-demo --wait=true --ignore-not-found >/dev/null
kubectl apply -f - <<'YAML' >/dev/null
apiVersion: v1
kind: Pod
metadata:
  name: qos-demo
  labels:
    app: qos-demo
spec:
  restartPolicy: Always
  containers:
    - name: app
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
      resources:
        requests:
          cpu: "100m"
          memory: "64Mi"
        limits:
          cpu: "100m"
          memory: "64Mi"
YAML
kubectl wait --for=condition=Ready pod/qos-demo --timeout=60s >/dev/null
