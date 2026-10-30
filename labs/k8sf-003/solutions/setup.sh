#!/usr/bin/env bash
set -euo pipefail

ip=$(kubectl get svc api -o jsonpath='{.spec.clusterIP}')

for i in $(seq 1 60); do
  kubectl get pods -n kube-system -l k8s-app=kube-dns --no-headers 2>/dev/null | grep -q . && break
  sleep 1
done
kubectl wait --for=condition=Ready pod -l k8s-app=kube-dns -n kube-system --timeout=60s >/dev/null 2>&1 || true

kubectl delete pod probe --ignore-not-found --now >/dev/null 2>&1 || true
kubectl run probe --image=busybox:1.36 --restart=Never --command -- \
  sh -c 'for i in $(seq 1 15); do wget -qO- -T 3 http://api.default.svc.cluster.local:80/ && exit 0; sleep 2; done; exit 1' >/dev/null

phase=""
for i in $(seq 1 90); do
  phase=$(kubectl get pod probe -o jsonpath='{.status.phase}' 2>/dev/null || true)
  if [ "$phase" = "Succeeded" ] || [ "$phase" = "Failed" ]; then
    break
  fi
  sleep 1
done

dns_worked="no"
if [ "$phase" = "Succeeded" ] && kubectl logs probe 2>/dev/null | grep -q "ok"; then
  dns_worked="yes"
fi
kubectl delete pod probe --ignore-not-found --now >/dev/null 2>&1 || true

cat <<EOF > /workspace/answers.env
SERVICE_CLUSTERIP=${ip}
DNS_LOOKUP_WORKED=${dns_worked}
EOF
