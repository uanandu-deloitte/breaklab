#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
# Oops: this manifest was copied from a different cluster with a different StorageClass name
kubectl apply -f - <<'YAML' >/dev/null
apiVersion: v1
kind: PersistentVolumeClaim
metadata: {name: data-pending}
spec:
  accessModes: [ReadWriteOnce]
  storageClassName: class-does-not-exist
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
sleep 5
