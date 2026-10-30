# Solution: Persistent Volumes and Claims

```bash
kubectl delete pod notes-writer --wait=true --ignore-not-found
kubectl delete pvc notes-pvc --wait=true --ignore-not-found

kubectl apply -f - <<'YAML'
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

kubectl wait --for=jsonpath='{.status.phase}'=Bound pvc/notes-pvc --timeout=30s
kubectl wait --for=condition=Ready pod/notes-writer --timeout=30s

kubectl exec notes-writer -- sh -c 'echo hello > /notes/testfile'
kubectl exec notes-writer -- cat /notes/testfile   # hello
```

The original PVC asked for `storageClassName: fast-ssd`, a class that
doesn't exist in this cluster at all, so nothing ever provisioned
against it. `storageClassName` is immutable once a PVC is created, so
the fix deletes the Pod and the PVC together and recreates both,
pointing the new PVC at `local-path` — k3s's built-in dynamic
provisioner, which supports the `ReadWriteOnce` access mode this
claim already correctly requested. Once the PVC binds, the Pod that
was waiting on it schedules and mounts the volume. Writing a file
through the mount and reading it back is what actually proves the
volume is usable — `Bound` alone only proves the control plane
matched the claim to provisioned storage, not that a Pod can read and
write through it.
