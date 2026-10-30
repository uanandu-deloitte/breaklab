#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

live_ip=$(kubectl get svc api -o jsonpath='{.spec.clusterIP}' 2>/dev/null)
[ -n "$live_ip" ] || fail "api Service not found"

[ -n "${SERVICE_CLUSTERIP:-}" ] || fail "SERVICE_CLUSTERIP is empty"
[ "$SERVICE_CLUSTERIP" = "$live_ip" ] || fail "recorded SERVICE_CLUSTERIP ($SERVICE_CLUSTERIP) does not match the Service's live ClusterIP ($live_ip)"

[ "${DNS_LOOKUP_WORKED:-}" = "yes" ] || fail "DNS_LOOKUP_WORKED must be recorded as yes"

# Independent proof: check.sh runs its own probe rather than trusting the
# learner's recorded answer for the functional claim.
for i in $(seq 1 60); do
  kubectl get pods -n kube-system -l k8s-app=kube-dns --no-headers 2>/dev/null | grep -q . && break
  sleep 1
done
kubectl wait --for=condition=Ready pod -l k8s-app=kube-dns -n kube-system --timeout=60s >/dev/null 2>&1 || true
kubectl delete pod check-dns-probe --ignore-not-found --now >/dev/null 2>&1 || true
kubectl run check-dns-probe --image=busybox:1.36 --restart=Never --command -- \
  sh -c 'for i in $(seq 1 15); do wget -qO- -T 3 http://api.default.svc.cluster.local:80/ && exit 0; sleep 2; done; exit 1' >/dev/null 2>&1

phase=""
for i in $(seq 1 90); do
  phase=$(kubectl get pod check-dns-probe -o jsonpath='{.status.phase}' 2>/dev/null)
  if [ "$phase" = "Succeeded" ] || [ "$phase" = "Failed" ]; then
    break
  fi
  sleep 1
done

output=$(kubectl logs check-dns-probe 2>/dev/null)
kubectl delete pod check-dns-probe --ignore-not-found --now >/dev/null 2>&1 || true

[ "$phase" = "Succeeded" ] || fail "independent DNS probe pod did not succeed (phase=$phase): $output"
echo "$output" | grep -q "ok" || fail "independent DNS probe did not get the expected response from api.default.svc.cluster.local"

echo "PASS: recorded Service ClusterIP matches live state, and DNS-based access to the api Service works"
