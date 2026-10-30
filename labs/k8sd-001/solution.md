# Solution: ConfigMaps, Two Ways

```bash
kubectl apply -f - <<'YAML'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: configmap-demo
  labels:
    app: configmap-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: configmap-demo
  template:
    metadata:
      labels:
        app: configmap-demo
    spec:
      containers:
        - name: app
          image: busybox:1.36
          command: ["sh", "-c", "sleep 3600"]
          env:
            - name: APP_LOG_LEVEL
              valueFrom:
                configMapKeyRef:
                  name: app-config
                  key: LOG_LEVEL
          volumeMounts:
            - name: config-volume
              mountPath: /etc/app
      volumes:
        - name: config-volume
          configMap:
            name: app-config
            items:
              - key: app.conf
                path: app.conf
YAML
kubectl rollout status deployment/configmap-demo --timeout=60s

kubectl exec deploy/configmap-demo -- printenv APP_LOG_LEVEL   # info
kubectl exec deploy/configmap-demo -- cat /etc/app/app.conf    # mode=production / retries=3
```

The container's `env[]` entry sources `APP_LOG_LEVEL` from the
ConfigMap's `LOG_LEVEL` key via `valueFrom.configMapKeyRef` — the
variable name and the key name don't have to match. The `app.conf` key
is handled differently: a `configMap` volume is declared at the Pod
level, restricted with `items` to just that one key so only
`/etc/app/app.conf` appears (not a second file for `LOG_LEVEL`), and
mounted into the container at `/etc/app`. Because the Deployment's Pod
template changed, Kubernetes rolls a new Pod that starts with both
pieces of wiring already in place — no live `kubectl exec` surgery
needed, which is also safer since there's no way to inject an
environment variable into an already-running process from the outside.
