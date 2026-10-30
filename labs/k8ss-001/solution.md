# Solution: Least Privilege with RBAC

```bash
# Read the live Role — this is the actual bug, not what the name implies.
kubectl get role pod-reader-role -n default -o yaml
# rules:
# - apiGroups: [""]
#   resources: ["pods"]
#   verbs: ["get", "delete"]        # missing "list", and "delete" shouldn't be there

# Confirm against the API server's own check, before touching anything.
kubectl auth can-i list pods   --as=system:serviceaccount:default:pod-reader-sa -n default   # no  (should be yes)
kubectl auth can-i delete pods --as=system:serviceaccount:default:pod-reader-sa -n default   # yes (should be no)

# Fix the rules to match the intended matrix exactly: get + list on pods, nothing else.
kubectl apply -f - <<'YAML'
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader-role
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list"]
YAML

# Verify every corner of the matrix, both the allows and the denies.
kubectl auth can-i get pods     --as=system:serviceaccount:default:pod-reader-sa -n default   # yes
kubectl auth can-i list pods    --as=system:serviceaccount:default:pod-reader-sa -n default   # yes
kubectl auth can-i create pods  --as=system:serviceaccount:default:pod-reader-sa -n default   # no
kubectl auth can-i update pods  --as=system:serviceaccount:default:pod-reader-sa -n default   # no
kubectl auth can-i delete pods  --as=system:serviceaccount:default:pod-reader-sa -n default   # no
kubectl auth can-i get secrets  --as=system:serviceaccount:default:pod-reader-sa -n default   # no
```

The `RoleBinding` was never the problem — it already pointed
`pod-reader-sa` at `pod-reader-role` correctly. The bug lived entirely
in the Role's `rules`: it was simultaneously too narrow (no `list`,
which breaks anything that needs to enumerate pods rather than fetch
one by name) and too broad (`delete`, which a read-only identity should
never hold). Replacing the rule with exactly `["get", "list"]` on
`pods` closes both gaps at once. `kubectl auth can-i --as=...` is the
only verification that actually matters here — it runs the same
`SubjectAccessReview` the API server runs for a real request from that
ServiceAccount, so it reports the truth about live permissions, not
just whether some YAML file was edited.
