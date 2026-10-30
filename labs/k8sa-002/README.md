# Default Limits with LimitRange

The `limits-demo` namespace already exists, with one Pod in it,
`pod-before` — created with **no** `resources` block at all, and there's
no `LimitRange` in the namespace yet, so it stayed exactly as written.

## Tasks

1. Create a `LimitRange` named `default-limits` in `limits-demo` with
   these exact defaults for every container in the namespace:

   ```yaml
   apiVersion: v1
   kind: LimitRange
   metadata:
     name: default-limits
     namespace: limits-demo
   spec:
     limits:
       - type: Container
         defaultRequest:
           cpu: 100m
           memory: 64Mi
         default:
           cpu: 200m
           memory: 128Mi
   ```

   Apply it.

2. Create a **new** Pod in `limits-demo` (any name other than
   `pod-before`, e.g. `pod-after`) whose manifest declares **no**
   `resources` block at all:

   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: pod-after
     namespace: limits-demo
   spec:
     restartPolicy: Never
     containers:
       - name: sleeper
         image: busybox:1.36
         command: ["sh", "-c", "sleep 3600"]
   ```

   Apply it and wait for it to reach `Running`.

3. Confirm Kubernetes injected the LimitRange's defaults into this Pod's
   live spec — you never typed a `resources` block, so anything you see
   here came from admission, not from you:

   ```
   kubectl get pod pod-after -n limits-demo -o jsonpath='{.spec.containers[0].resources}'
   ```

4. Optional but worth doing: check `pod-before` the same way. It has no
   `resources` set, and it still won't after you add the LimitRange — a
   LimitRange only affects objects created *after* it exists.

## Useful commands

```
kubectl get limitrange -n limits-demo
kubectl describe limitrange default-limits -n limits-demo
kubectl get pod pod-after -n limits-demo -o jsonpath='{.spec.containers[0].resources}'
kubectl get pod pod-before -n limits-demo -o jsonpath='{.spec.containers[0].resources}'
```

Get a shell: `docker exec -it breaklab-k8sa-002 bash`
