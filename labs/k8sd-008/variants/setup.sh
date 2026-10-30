#!/usr/bin/env bash
set -euo pipefail
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
wait-for-k8s
helm install web /chart
kubectl rollout status deployment/web-web --timeout=60s
