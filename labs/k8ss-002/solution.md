# Solution: Run as Non-Root

```bash
# Confirm the starting state: no securityContext, running as root.
kubectl get pod scratch-writer -o jsonpath='{.spec.securityContext}{"\n"}'
kubectl exec scratch-writer -- id -u        # 0

# Edit the manifest that's baked into this image, adding:
#   - spec.securityContext.runAsNonRoot: true
#   - spec.securityContext.runAsUser: 10001
#   - containers[0].securityContext.readOnlyRootFilesystem: true
#   - a volume + volumeMount so /data stays writable
vi /manifests/pod.yaml
```

```yaml
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
```

```bash
# A Pod's securityContext is immutable once running — recreate it.
kubectl delete pod scratch-writer
kubectl apply -f /manifests/pod.yaml
kubectl wait --for=condition=Ready pod/scratch-writer --timeout=60s

# Verify every part of the fix independently.
kubectl exec scratch-writer -- id -u                                    # 10001
kubectl exec scratch-writer -- sh -c 'echo hi > /data/t && cat /data/t'  # hi — still writable
kubectl exec scratch-writer -- sh -c 'touch /nope'                      # Read-only file system
```

Nothing about `readOnlyRootFilesystem` understands "this app happens to
need `/data`" — it locks down the entire image layer indiscriminately.
The only reason the container keeps working is the `emptyDir` volume
explicitly mounted at `/data`, which isn't part of the image layer at
all and stays writable regardless of the rest of the filesystem's
state. Skip that volume and the fix "works" (the securityContext fields
are all set correctly) right up until the container's own write to
`/data/output.txt` starts failing — a mistake that's easy to miss if you
only check the Pod's spec and never confirm it's still actually doing
its job.
