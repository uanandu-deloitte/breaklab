# Namespaces and Resource Quotas

The `quota-demo` namespace already exists, with a `ResourceQuota` named
`pod-quota` capping it at **2 Pods total**, and one Pod (`pod-a`) already
running there — using 1 of those 2 slots.

## Tasks

1. Create a second Pod in `quota-demo` so the namespace reaches its quota
   limit. Save this manifest (e.g. as `/workspace/pod-b.yaml`) and apply it:

   ```yaml
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
   ```

   Confirm it reaches `Running` before moving on.

2. Confirm the namespace is now at its limit:

   ```
   kubectl get resourcequota pod-quota -n quota-demo
   ```

   `pods` should show `2/2` used/hard.

3. Now try to create a **third** Pod in the same namespace (any name, same
   shape as above works fine — call it `pod-c`). Do not expect this to
   succeed — the point is to watch the API server reject it and read its
   real error text.

4. Record the exact rejection reason you saw — or a recognizable piece of
   it — in `/workspace/answers.env`:

   ```
   REJECTION_REASON="<the real error message kubectl printed>"
   ```

   Don't guess the wording. Run the apply, copy what it actually says.

## Useful commands

```
kubectl get resourcequota -n quota-demo
kubectl get pods -n quota-demo
kubectl apply -f /workspace/pod-b.yaml
kubectl apply -f /workspace/pod-c.yaml   # expected to fail
```

Get a shell: `docker exec -it breaklab-k8sa-001 bash`
