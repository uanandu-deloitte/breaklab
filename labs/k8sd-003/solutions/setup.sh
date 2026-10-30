#!/usr/bin/env bash
set -euo pipefail
kubectl delete pod notes-writer --wait=true --ignore-not-found >/dev/null
kubectl delete pvc notes-pvc --wait=true --ignore-not-found >/dev/null
kubectl apply -f - <<'YAML' >/dev/null
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: notes-pvc
spec:
  accessModes: [ReadWriteOnce]
  storageClassName: local-path
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: Pod
metadata:
  name: notes-writer
  labels:
    app: notes-writer
spec:
  restartPolicy: Never
  containers:
    - name: writer
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
      volumeMounts:
        - name: notes
          mountPath: /notes
  volumes:
    - name: notes
      persistentVolumeClaim:
        claimName: notes-pvc
YAML
kubectl wait --for=jsonpath='{.status.phase}'=Bound pvc/notes-pvc --timeout=30s >/dev/null
kubectl wait --for=condition=Ready pod/notes-writer --timeout=30s >/dev/null
