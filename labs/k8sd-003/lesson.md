# What actually has to be true for a PVC to bind

A PersistentVolumeClaim is a request, not a guarantee. It only turns
into a bound, usable volume when every one of these lines up:

- **The StorageClass it names has to exist.** `storageClassName:
  fast-ssd` binds to nothing if no StorageClass called `fast-ssd`
  exists in this cluster — it just sits `Pending` forever, with no
  error anywhere except the PVC's own Events.
- **That StorageClass's provisioner has to support what's being
  asked for.** Capacity, and — critically — access mode. A
  provisioner that only implements `ReadWriteOnce` (most node-local
  ones do) will never satisfy a claim that asks for `ReadWriteMany`,
  no matter how long it waits.
- **The claim has to actually get looked at.** A Pod referencing a
  PVC in its `volumes[]` just waits patiently for the PVC to bind —
  it's not an error state on the Pod's side, which is why the real
  diagnosis always starts at the PVC, not the Pod.

k3s ships one built-in StorageClass out of the box: `local-path`,
backed by the Rancher local-path-provisioner. It dynamically
provisions a directory on the node's own disk per claim and supports
`ReadWriteOnce` — which is exactly what a single Pod mounting its own
private volume needs.

## Why "Bound" isn't the same as "working"

A PVC reaching `Bound` only proves the control plane successfully
matched it to provisioned storage — it says nothing about whether a
Pod can actually read and write through the resulting mount. A
misconfigured provisioner, a network filesystem that's bound but
unreachable from this particular node, or wrong permissions on the
provisioned directory can all produce a `Bound` PVC that's still
useless to the Pod trying to use it. The only way to be sure is the
same way you'd prove any filesystem works: write something through
the mount, then read it back, from inside the Pod that's actually
using it.
