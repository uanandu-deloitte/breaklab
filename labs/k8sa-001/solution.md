# Solution: Namespaces and Resource Quotas

```bash
# 1. Bring the namespace to its quota limit.
cat <<'YAML' | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: pod-b
  namespace: quota-demo
spec:
  restartPolicy: Never
  containers:
    - name: sleeper
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
YAML
kubectl wait --for=condition=Ready pod/pod-b -n quota-demo --timeout=60s

# 2. Confirm the quota is now fully used.
kubectl get resourcequota pod-quota -n quota-demo
# pods: 2/2

# 3. Try a third Pod — the API server rejects it immediately.
cat <<'YAML' | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: pod-c
  namespace: quota-demo
spec:
  restartPolicy: Never
  containers:
    - name: sleeper
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
YAML
# Error from server (Forbidden): error when creating "STDIN": pods "pod-c"
# is forbidden: exceeded quota: pod-quota, requested: pods=1, used: pods=2,
# limited: pods=2

# 4. Record the real rejection text.
echo 'REJECTION_REASON="pods \"pod-c\" is forbidden: exceeded quota: pod-quota, requested: pods=1, used: pods=2, limited: pods=2"' >> /workspace/answers.env
```

`pod-c` was never created — the quota admission plugin rejects it before
it's ever persisted, so `pod-c` never shows up in `kubectl get pods`
even transiently. The namespace stays at exactly 2 Pods (`pod-a` and
`pod-b`), both Running, and the quota stays at `2/2` forever until one of
them is deleted.
