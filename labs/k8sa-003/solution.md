# Solution: Taints, Tolerations, and Scheduling

```bash
# 1. Confirm the diagnosis.
kubectl get pod tainted-workload -o wide
# tainted-workload   0/1   Pending   0   ...   <none>   <none>

kubectl describe pod tainted-workload
# Events:
#   Warning  FailedScheduling  ...  0/1 nodes are available:
#   1 node(s) had untolerated taint {dedicated: special}.

node=$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')
kubectl describe node "$node" | grep -A1 Taints
# Taints: dedicated=special:NoSchedule

# 2. Tolerations can't be patched onto an existing Pod — delete and recreate.
kubectl delete pod tainted-workload --wait=true

cat <<'YAML' | kubectl apply -f -
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
YAML
kubectl wait --for=condition=Ready pod/tainted-workload --timeout=60s

# 3. Confirm it's scheduled and Running, and the node's taint is untouched.
kubectl get pod tainted-workload -o wide
kubectl describe node "$node" | grep -A1 Taints
```

The taint never moved — it's still on the node exactly as it was. What
changed is the Pod: its new spec declares a toleration that matches the
taint's key, value, and effect precisely, which is the one thing that
lets the scheduler consider this node a valid placement for it.
