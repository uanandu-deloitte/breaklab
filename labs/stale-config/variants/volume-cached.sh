#!/usr/bin/env bash
set -euo pipefail
wait-for-k8s
kubectl apply -f /manifests/configmap.yaml >/dev/null
kubectl apply -f /manifests/deploy-volume.yaml >/dev/null
kubectl wait --for=condition=Available deployment/configwatcher --timeout=60s >/dev/null
# Oops: ops just flipped the ConfigMap off, but the app only reads its mounted
# config file once at startup, not on every request
kubectl patch configmap site-config --type merge -p '{"data":{"mode":"normal"}}' >/dev/null
