# What has to exist for a PVC to bind

A PersistentVolumeClaim's `storageClassName` is a request, not a
guarantee. For a PVC to move from `Pending` to `Bound`, a
`StorageClass` with that exact name has to actually exist in the
cluster, and its provisioner has to be capable of satisfying what the
claim asks for (capacity, access mode, and so on). Miss either
condition and the PVC just sits `Pending` — indefinitely, with no
timeout and no escalation.

## The simplest way this breaks

The single most common cause is also the quietest: a manifest copied
from a different cluster, a different example, or an older version of
the same cluster, referencing a `storageClassName` that was real
*somewhere*, but isn't real *here*. The YAML is completely
well-formed. `kubectl apply` admits it without complaint — there's
nothing about the shape of the object that's wrong, only its content.
Nothing provisions against it, because nothing matching that name
exists, and the PVC just waits forever.

## Where the evidence actually lives

- `kubectl get pvc <name>` shows the phase (`Pending`), but not why.
- `kubectl describe pvc <name>`'s Events section at the bottom is
  where the provisioning controller actually records its rejection —
  this is the first thing to read, not the last resort.
- `kubectl get pvc <name> -o jsonpath='{.spec.storageClassName}'`
  shows exactly what the PVC is requesting, with no ambiguity.
- `kubectl get storageclass` shows every StorageClass that genuinely
  exists right now. k3s ships exactly one out of the box: `local-path`.
  If the PVC's requested name isn't in this list, character for
  character, it will never bind — no amount of waiting fixes it.

## Why the Pod waiting on it looks fine

A Pod that mounts a PVC in its `volumes[]` just waits for that PVC to
bind — from the Pod's own perspective, this isn't an error, it's a
normal (if prolonged) pending state. Nothing shows up in the Pod's own
events or logs explaining the real cause, which is exactly why the
diagnosis always has to start at the PVC itself, never at the Pod
that's merely waiting on it.
