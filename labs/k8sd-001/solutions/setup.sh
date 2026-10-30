#!/usr/bin/env bash
set -euo pipefail
kubectl apply -f - <<'YAML' >/dev/null
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
kubectl rollout status deployment/configmap-demo --timeout=60s >/dev/null
