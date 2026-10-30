# StorageClasses: What Has to Match for a PVC to Bind

The `reports-pvc` PersistentVolumeClaim has sat `Pending` since it was
created, and `reports-writer` — the Pod waiting to mount it — is stuck
waiting right along with it. Diagnose exactly why before you'd reach
for a fix.

## Tasks

1. `kubectl get pvc reports-pvc` — confirm the phase really is
   `Pending`.
2. `kubectl describe pvc reports-pvc` — read the Events at the bottom
   for the provisioner's own rejection reason.
3. `kubectl get pvc reports-pvc -o jsonpath='{.spec.storageClassName}'`
   — see exactly what the PVC is requesting.
4. `kubectl get storageclass` — see what StorageClasses actually
   exist in this cluster.
5. `kubectl get pod -l app=reports-writer` — confirm the waiting Pod
   itself shows a normal pending state, not an error.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
PVC_PHASE=
REQUESTED_STORAGECLASS=
AVAILABLE_STORAGECLASSES=
CLASS_EXISTS=
POD_PHASE=
```

- `PVC_PHASE`: `reports-pvc`'s current phase.
- `REQUESTED_STORAGECLASS`: the exact `storageClassName` the PVC
  requests.
- `AVAILABLE_STORAGECLASSES`: every StorageClass name that actually
  exists in this cluster, space-separated (if there's more than one).
- `CLASS_EXISTS`: `yes` or `no` — does `REQUESTED_STORAGECLASS` appear
  in `AVAILABLE_STORAGECLASSES`?
- `POD_PHASE`: `reports-writer`'s current phase.

Collect every value from the running cluster itself — don't hard-code
anything from the lesson.

Get a shell: `docker exec -it breaklab-k8sd-006 bash`
