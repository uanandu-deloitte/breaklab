#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

kubectl get deployment web-web >/dev/null 2>&1 || fail "deployment web-web doesn't exist -- was helm install ever run?"
kubectl rollout status deployment/web-web --timeout=30s >/dev/null 2>&1 || fail "deployment web-web isn't ready"

run_exit=1
out=""
for attempt in 1 2 3 4 5; do
  kubectl delete pod helm-check --ignore-not-found --wait=true >/dev/null 2>&1
  out=$(kubectl run helm-check --image=busybox:1.36 --restart=Never --rm -i --quiet \
    --command -- wget -qO- --timeout=5 http://web-web:8080 2>&1)
  run_exit=$?
  [ "$run_exit" -eq 0 ] && break
  sleep 2
done

[ "$run_exit" -eq 0 ] || fail "wget against the real Service (web-web:8080) failed -- the rendered Service still doesn't reach the container's real port:
$out"

echo "$out" | grep -qi "Server address" || fail "got a response, but not from the real app: $out"

echo "PASS: the real Service routes real traffic to the container's actual listening port"
