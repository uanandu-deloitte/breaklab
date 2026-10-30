# Least Privilege with RBAC

This namespace (`default`) has a `ServiceAccount` named `pod-reader-sa`, a
`Role` named `pod-reader-role`, and a `RoleBinding` that already connects
them. The intent was simple: **`pod-reader-sa` should be able to `get`
and `list` pods in this namespace, and nothing else** — no `create`, no
`update`, no `delete`, and no access to any other resource type (for
example, `secrets`).

The Role currently in the cluster does not match that. Read its actual
rules and fix them.

## Tasks

- inspect the live `Role`'s rules (don't assume — read them)
- fix `pod-reader-role` so its rules grant **exactly**: `get` and `list`
  on `pods`, and nothing else
- verify the fix yourself with `kubectl auth can-i <verb> <resource>
  --as=system:serviceaccount:default:pod-reader-sa -n default` for each
  verb that matters: `get`, `list`, `create`, `update`, `delete` (on
  `pods`), and `get` on `secrets`

## Useful commands

```
kubectl get role pod-reader-role -n default -o yaml
kubectl get rolebinding pod-reader-binding -n default -o yaml
kubectl auth can-i get pods --as=system:serviceaccount:default:pod-reader-sa -n default
kubectl auth can-i list pods --as=system:serviceaccount:default:pod-reader-sa -n default
kubectl auth can-i delete pods --as=system:serviceaccount:default:pod-reader-sa -n default
kubectl auth can-i get secrets --as=system:serviceaccount:default:pod-reader-sa -n default
```

Get a shell: `docker exec -it breaklab-k8ss-001 bash`
