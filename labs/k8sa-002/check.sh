#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

ns=limits-demo
lr_name=default-limits

req_cpu=$(kubectl get limitrange "$lr_name" -n "$ns" -o jsonpath='{.spec.limits[0].defaultRequest.cpu}' 2>/dev/null)
req_mem=$(kubectl get limitrange "$lr_name" -n "$ns" -o jsonpath='{.spec.limits[0].defaultRequest.memory}' 2>/dev/null)
lim_cpu=$(kubectl get limitrange "$lr_name" -n "$ns" -o jsonpath='{.spec.limits[0].default.cpu}' 2>/dev/null)
lim_mem=$(kubectl get limitrange "$lr_name" -n "$ns" -o jsonpath='{.spec.limits[0].default.memory}' 2>/dev/null)

[ -n "$req_cpu" ] || fail "LimitRange $lr_name not found in namespace $ns (or has no defaultRequest.cpu)"
[ "$req_cpu" = "100m" ] || fail "LimitRange defaultRequest.cpu is ${req_cpu:-unset}, expected 100m"
[ "$req_mem" = "64Mi" ] || fail "LimitRange defaultRequest.memory is ${req_mem:-unset}, expected 64Mi"
[ "$lim_cpu" = "200m" ] || fail "LimitRange default.cpu is ${lim_cpu:-unset}, expected 200m"
[ "$lim_mem" = "128Mi" ] || fail "LimitRange default.memory is ${lim_mem:-unset}, expected 128Mi"

new_pod=$(kubectl get pods -n "$ns" -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' 2>/dev/null | grep -v '^pod-before$' | head -1)
[ -n "$new_pod" ] || fail "no new pod found in $ns besides pod-before"

phase=$(kubectl get pod "$new_pod" -n "$ns" -o jsonpath='{.status.phase}' 2>/dev/null)
[ "$phase" = "Running" ] || fail "$new_pod is ${phase:-missing}, not Running"

pod_req_cpu=$(kubectl get pod "$new_pod" -n "$ns" -o jsonpath='{.spec.containers[0].resources.requests.cpu}' 2>/dev/null)
pod_req_mem=$(kubectl get pod "$new_pod" -n "$ns" -o jsonpath='{.spec.containers[0].resources.requests.memory}' 2>/dev/null)
pod_lim_cpu=$(kubectl get pod "$new_pod" -n "$ns" -o jsonpath='{.spec.containers[0].resources.limits.cpu}' 2>/dev/null)
pod_lim_mem=$(kubectl get pod "$new_pod" -n "$ns" -o jsonpath='{.spec.containers[0].resources.limits.memory}' 2>/dev/null)

[ "$pod_req_cpu" = "$req_cpu" ] || fail "$new_pod requests.cpu is ${pod_req_cpu:-unset}, expected the LimitRange default $req_cpu"
[ "$pod_req_mem" = "$req_mem" ] || fail "$new_pod requests.memory is ${pod_req_mem:-unset}, expected the LimitRange default $req_mem"
[ "$pod_lim_cpu" = "$lim_cpu" ] || fail "$new_pod limits.cpu is ${pod_lim_cpu:-unset}, expected the LimitRange default $lim_cpu"
[ "$pod_lim_mem" = "$lim_mem" ] || fail "$new_pod limits.memory is ${pod_lim_mem:-unset}, expected the LimitRange default $lim_mem"

echo "PASS: $new_pod picked up the namespace's LimitRange defaults automatically"
