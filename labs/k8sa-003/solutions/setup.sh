#!/usr/bin/env bash
set -euo pipefail
kubectl delete pod tainted-workload --wait=true
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: Pod
metadata:
  name: tainted-workload
  labels:
    app: tainted-workload
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
