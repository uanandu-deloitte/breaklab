#!/usr/bin/env bash
set -euo pipefail
kubectl delete pod scratch-writer --ignore-not-found --wait=true >/dev/null
cat <<'YAML' | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: scratch-writer
  namespace: default
  labels:
    app: scratch-writer
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 10001
  containers:
    - name: scratch-writer
      image: busybox:1.36
      command: ["sh", "-c", "mkdir -p /data; while true; do date >> /data/output.txt; sleep 2; done"]
      securityContext:
        readOnlyRootFilesystem: true
      volumeMounts:
        - name: scratch
          mountPath: /data
  volumes:
    - name: scratch
      emptyDir: {}
YAML
kubectl wait --for=condition=Ready pod/scratch-writer --timeout=60s >/dev/null
