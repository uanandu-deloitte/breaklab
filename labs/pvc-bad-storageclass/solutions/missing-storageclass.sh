#!/usr/bin/env bash
set -euo pipefail
kubectl delete pod data-consumer --wait=true --ignore-not-found
kubectl delete pvc data-pending --wait=true --ignore-not-found
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: PersistentVolumeClaim
metadata: {name: data-pending}
spec:
  accessModes: [ReadWriteOnce]
  storageClassName: local-path
  resources: {requests: {storage: 1Gi}}
---
apiVersion: v1
kind: Pod
metadata: {name: data-consumer}
spec:
  restartPolicy: Never
  containers:
  - name: sleeper
    image: busybox:1.36
    command: ["sh","-c","sleep 3600"]
    volumeMounts:
    - name: data
      mountPath: /data
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: data-pending
YAML
kubectl wait --for=jsonpath='{.status.phase}'=Bound pvc/data-pending --timeout=30s
