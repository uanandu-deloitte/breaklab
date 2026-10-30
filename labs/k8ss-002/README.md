# Run as Non-Root

There's a Pod named `scratch-writer` already running in this cluster. It
has **no `securityContext` at all** — it runs as root, with a fully
read-write root filesystem, and it works. That's exactly the problem:
nothing here would stop it from writing anywhere on its own filesystem,
or behaving as root if something inside it were ever compromised.

The container's job is simple: it appends a timestamp to
`/data/output.txt` every two seconds, forever.

## Tasks

Harden the Pod's `securityContext` so that:

- `runAsNonRoot: true`
- `runAsUser: 10001` (exactly this uid)
- `readOnlyRootFilesystem: true` on the container

...while keeping the container able to keep doing its actual job: writing
to `/data/output.txt`. Since the root filesystem becomes read-only, that
means adding a volume — an `emptyDir` mounted at `/data` is enough.

A running Pod's `securityContext` can't be edited in place (it's part of
an immutable field). The manifest that was applied at startup is sitting
at `/manifests/pod.yaml` inside this container — edit it, then:

```
kubectl delete pod scratch-writer
kubectl apply -f /manifests/pod.yaml
```

## Verify it yourself

```
kubectl get pod scratch-writer -o yaml | grep -A5 securityContext
kubectl exec scratch-writer -- id -u                       # should print 10001
kubectl exec scratch-writer -- sh -c 'echo hi > /data/t && cat /data/t'   # should work
kubectl exec scratch-writer -- sh -c 'touch /nope'          # should fail: read-only file system
```

Get a shell: `docker exec -it breaklab-k8ss-002 bash`
