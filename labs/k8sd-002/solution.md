# Solution: Requests, Limits, and QoS

```bash
kubectl delete pod qos-demo --wait=true
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: Pod
metadata:
  name: qos-demo
  labels:
    app: qos-demo
spec:
  restartPolicy: Always
  containers:
    - name: app
      image: busybox:1.36
      command: ["sh", "-c", "sleep 3600"]
      resources:
        requests:
          cpu: "100m"
          memory: "64Mi"
        limits:
          cpu: "100m"
          memory: "64Mi"
YAML
kubectl wait --for=condition=Ready pod/qos-demo --timeout=60s

kubectl get pod qos-demo -o jsonpath='{.status.qosClass}'   # Guaranteed
```

Setting `requests.cpu` and `limits.cpu` to the same `100m`, and
`requests.memory` and `limits.memory` to the same `64Mi`, on the Pod's
only container is exactly what Guaranteed requires: both resources
set on every container, request equal to limit on both. Because a
Pod's resources can't be patched in place, the fix deletes the old Pod
and applies a new one with the corrected spec rather than trying to
edit the running one. `64Mi` is comfortably more than an idle
`sleep`-looping busybox actually needs, so the Pod stays Running
instead of trading its new QoS class for an immediate OOM kill.
