# Taints, Tolerations, and Scheduling

The cluster's one node has been tainted for a dedicated workload pool:

```
dedicated=special:NoSchedule
```

A Pod named `tainted-workload` was applied with no toleration for that
taint, so it's sitting `Pending` — there's only one node in this
cluster, and the scheduler won't place an untolerated Pod on a tainted
one, no matter how idle it is.

## Tasks

1. Confirm the diagnosis yourself:

   ```
   kubectl get pod tainted-workload -o wide
   kubectl describe pod tainted-workload
   kubectl describe node $(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')
   ```

   Look for `node(s) had untolerated taint` in the Pod's Events, and the
   exact `dedicated=special:NoSchedule` taint on the node.

2. A Pod's `tolerations` are set at creation and can't be patched onto a
   running/pending Pod's spec. Delete `tainted-workload` and recreate it
   with a toleration that matches the node's taint exactly — same key,
   same value, same effect:

   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: tainted-workload
   spec:
     restartPolicy: Never
     tolerations:
       - key: dedicated
         operator: Equal
         value: special
         effect: NoSchedule
     containers:
       - name: sleeper
         image: busybox:1.36
         command: ["sh", "-c", "sleep 3600"]
   ```

3. Do **not** remove the node's taint to "fix" this — the check verifies
   the taint is still there. The fix belongs on the Pod's side.

4. Confirm `tainted-workload` reaches `Running`.

## Useful commands

```
kubectl get pod tainted-workload -o wide
kubectl describe pod tainted-workload
kubectl describe node <node-name>
kubectl delete pod tainted-workload --wait=true
kubectl get pod tainted-workload -o jsonpath='{.spec.tolerations}'
```

Get a shell: `docker exec -it breaklab-k8sa-003 bash`
