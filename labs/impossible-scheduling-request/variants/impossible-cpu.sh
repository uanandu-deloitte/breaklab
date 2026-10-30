#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
# Oops: a template for a much bigger node pool was used here by mistake
kubectl apply -f - <<'YAML' >/dev/null
apiVersion: v1
kind: Pod
metadata: {name: worker-task}
spec:
  restartPolicy: Never
  containers:
  - name: sleeper
    image: busybox:1.36
    command: ["sh","-c","sleep 3600"]
    resources:
      requests: {cpu: "1000", memory: 16Mi}
YAML
sleep 3
