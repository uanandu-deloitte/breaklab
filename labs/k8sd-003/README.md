# Persistent Volumes and Claims

The `notes-pvc` PersistentVolumeClaim has been `Pending` since it was
created, and `notes-writer` — the Pod that wants to mount it — is
stuck waiting right along with it.

## Tasks

1. `kubectl get pvc notes-pvc` — confirm the phase really is
   `Pending`.
2. `kubectl describe pvc notes-pvc` — read the Events at the bottom
   for the provisioner's own reason.
3. `kubectl get storageclass` — see what StorageClasses actually
   exist in this cluster (k3s ships exactly one, `local-path`, by
   default).
4. `kubectl get pvc notes-pvc -o jsonpath='{.spec.storageClassName}'`
   — compare what the PVC is asking for against what's actually
   there.
5. Fix the PVC's `storageClassName` and re-apply it (delete and
   recreate the PVC and the Pod together, since a PVC's
   `storageClassName` is immutable once created).
6. Confirm `kubectl get pvc notes-pvc` shows `Bound` and
   `kubectl get pod notes-writer` shows `Running`.
7. Prove it's actually usable, not just bound: exec in and write,
   then read back, a file under `/notes`.

Get a shell: `docker exec -it breaklab-k8sd-003 bash`
