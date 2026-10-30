#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

node=$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
[ -n "$node" ] || fail "no node found"

tk=$(kubectl get node "$node" -o jsonpath='{.spec.taints[0].key}' 2>/dev/null)
tv=$(kubectl get node "$node" -o jsonpath='{.spec.taints[0].value}' 2>/dev/null)
te=$(kubectl get node "$node" -o jsonpath='{.spec.taints[0].effect}' 2>/dev/null)
if [ "$tk" != "dedicated" ] || [ "$tv" != "special" ] || [ "$te" != "NoSchedule" ]; then
  fail "node $node no longer carries the dedicated=special:NoSchedule taint (key=${tk:-none} value=${tv:-none} effect=${te:-none}) — add a toleration instead of removing the taint"
fi

phase=$(kubectl get pod tainted-workload -o jsonpath='{.status.phase}' 2>/dev/null)
[ "$phase" = "Running" ] || fail "tainted-workload is ${phase:-missing}, not Running"

match=$(kubectl get pod tainted-workload -o jsonpath='{range .spec.tolerations[*]}{.key}={.value}:{.effect}{"\n"}{end}' 2>/dev/null | grep -x "dedicated=special:NoSchedule" || true)
[ -n "$match" ] || fail "tainted-workload has no toleration matching the node's dedicated=special:NoSchedule taint"

echo "PASS: tainted-workload tolerates the node's taint and is Running"
