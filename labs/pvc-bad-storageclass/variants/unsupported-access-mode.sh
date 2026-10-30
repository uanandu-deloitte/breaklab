#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
# Oops: this workload actually needs ReadWriteOnce, but was templated from
# a shared-storage example that requested ReadWriteMany, which this
# cluster's provisioner doesn't support
kubectl apply -f - <<'YAML' >/dev/null
apiVersion: v1
kind: PersistentVolumeClaim
metadata: {name: data-pending}
spec:
  accessModes: [ReadWriteMany]
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
sleep 5
