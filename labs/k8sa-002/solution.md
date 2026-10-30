# Solution: Default Limits with LimitRange

```bash
# 1. Create the LimitRange with the required defaults.
cat <<'YAML' | kubectl apply -f -
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
YAML

# 2. Create a new Pod with NO resources block of its own.
cat <<'YAML' | kubectl apply -f -
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
YAML
kubectl wait --for=condition=Ready pod/pod-after -n limits-demo --timeout=60s

# 3. Confirm the defaults were injected.
kubectl get pod pod-after -n limits-demo -o jsonpath='{.spec.containers[0].resources}'
# {"limits":{"cpu":"200m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"64Mi"}}

# 4. Confirm pod-before was unaffected (created before the LimitRange existed).
kubectl get pod pod-before -n limits-demo -o jsonpath='{.spec.containers[0].resources}'
# {}
```

`pod-after`'s manifest never mentioned CPU or memory — the
`{"limits":...,"requests":...}` block came entirely from the
`default-limits` LimitRange, injected the instant the Pod was admitted.
`pod-before` proves the injection is admission-time only: it predates the
LimitRange, so it still has an empty `resources` object, and nothing will
ever retroactively patch it.
