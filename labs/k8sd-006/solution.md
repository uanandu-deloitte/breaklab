# Solution: StorageClasses: What Has to Match for a PVC to Bind

```bash
phase=$(kubectl get pvc reports-pvc -o jsonpath='{.status.phase}')
requested=$(kubectl get pvc reports-pvc -o jsonpath='{.spec.storageClassName}')
available=$(kubectl get storageclass -o jsonpath='{.items[*].metadata.name}')
pod_phase=$(kubectl get pod reports-writer -o jsonpath='{.status.phase}')

cat <<EOF > /workspace/answers.env
PVC_PHASE=$phase
REQUESTED_STORAGECLASS=$requested
AVAILABLE_STORAGECLASSES=$available
CLASS_EXISTS=no
POD_PHASE=$pod_phase
EOF
```

`kubectl describe pvc reports-pvc`'s Events show the provisioning
controller's own rejection: no StorageClass named `fast-ssd` exists in
this cluster. `kubectl get storageclass` confirms it directly — k3s
ships exactly one built-in class, `local-path`, and `fast-ssd` simply
isn't in that list. The claim's YAML is entirely well-formed; nothing
about its shape is wrong, only its content — this is exactly the kind
of mismatch that slips through code review because `kubectl apply`
never complains about it.

`reports-writer`, the Pod mounting this PVC, shows `Pending` too, but
for an entirely uninteresting reason: it's just waiting on a volume
that will never materialize. Nothing in the Pod's own events explains
why — that evidence lives only in the PVC's own Events, which is why
the diagnosis has to start there, not with the Pod that's merely
waiting.
