#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

action=$(kubectl get clusterpolicy disallow-privileged-containers -o jsonpath='{.spec.validationFailureAction}' 2>/dev/null)
[ "$action" = "Enforce" ] || fail "disallow-privileged-containers' validationFailureAction is \"${action:-missing}\", expected \"Enforce\""

if kubectl get pod bad-pod >/dev/null 2>&1; then
  fail "bad-pod (a real privileged Pod that violates the policy) is still running -- Enforce only blocks future admissions, it doesn't remove what already slipped through"
fi

out=$(kubectl apply -f /manifests/bad-pod.yaml 2>&1)
apply_exit=$?
if [ "$apply_exit" -eq 0 ]; then
  kubectl delete pod bad-pod --wait=true >/dev/null 2>&1
  fail "a privileged Pod was admitted -- the policy still isn't actually blocking anything"
fi

echo "$out" | grep -qi "privileged" || fail "pod creation failed, but not for the expected policy reason: $out"

echo "PASS: disallow-privileged-containers is set to Enforce and a real privileged Pod is rejected by Kyverno's admission webhook"
