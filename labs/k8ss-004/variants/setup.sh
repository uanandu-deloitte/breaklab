#!/usr/bin/env bash
set -euo pipefail

wait-for-k8s

# --server-side avoids kubectl apply's last-applied-configuration
# annotation size limit, which the real Kyverno CRDs exceed.
kubectl apply --server-side -f /manifests/kyverno-install.yaml

kubectl -n kyverno rollout status deployment/kyverno-admission-controller --timeout=180s
kubectl -n kyverno rollout status deployment/kyverno-background-controller --timeout=180s
kubectl -n kyverno rollout status deployment/kyverno-reports-controller --timeout=180s
kubectl -n kyverno rollout status deployment/kyverno-cleanup-controller --timeout=180s

kubectl apply -f /manifests/policy.yaml
sleep 3
# BUG left for the learner to find: the policy is Audit-only, so this
# privileged Pod gets admitted even though it violates it.
kubectl apply -f /manifests/bad-pod.yaml
