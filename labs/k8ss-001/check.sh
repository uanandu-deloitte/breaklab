#!/usr/bin/env bash
set -uo pipefail

NS=default
SA="system:serviceaccount:${NS}:pod-reader-sa"

fail() { echo "FAIL: $1"; exit 1; }

can_i() {
  # kubectl auth can-i prints exactly "yes" or "no" to stdout; its exit
  # code mirrors that, so capture stdout and ignore the exit code here.
  kubectl auth can-i "$1" "$2" --as="$SA" -n "$NS" 2>/dev/null
}

result=$(can_i get pods) || true
[ "$result" = "yes" ] || fail "pod-reader-sa should be able to 'get' pods in $NS, but can-i says: $result"

result=$(can_i list pods) || true
[ "$result" = "yes" ] || fail "pod-reader-sa should be able to 'list' pods in $NS, but can-i says: $result"

result=$(can_i create pods) || true
[ "$result" = "no" ] || fail "pod-reader-sa should NOT be able to 'create' pods, but can-i says: $result"

result=$(can_i delete pods) || true
[ "$result" = "no" ] || fail "pod-reader-sa should NOT be able to 'delete' pods, but can-i says: $result"

result=$(can_i update pods) || true
[ "$result" = "no" ] || fail "pod-reader-sa should NOT be able to 'update' pods, but can-i says: $result"

result=$(can_i get secrets) || true
[ "$result" = "no" ] || fail "pod-reader-sa should NOT have any access to secrets, but can-i says: $result"

echo "PASS: pod-reader-sa's live RBAC matches the exact least-privilege matrix (get+list pods only)"
