#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

REQUIRED_REPLICAS=4

live_replicas=$(kubectl get deployment hello -o jsonpath='{.spec.replicas}' 2>/dev/null)
[ -n "$live_replicas" ] || fail "hello Deployment not found"
[ "$live_replicas" = "$REQUIRED_REPLICAS" ] || fail "deployment/hello has $live_replicas replicas, expected $REQUIRED_REPLICAS"

ready_replicas=$(kubectl get deployment hello -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
[ "$ready_replicas" = "$REQUIRED_REPLICAS" ] || fail "deployment/hello only has $ready_replicas/$REQUIRED_REPLICAS Pods Ready"

live_image=$(kubectl get deployment hello -o jsonpath='{.spec.template.spec.containers[0].image}' 2>/dev/null)
manifest_initial=$(awk -F': *' '/replicas:/{print $2; exit}' /manifests/app.yaml)

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${INITIAL_REPLICA_COUNT:-}" ] || fail "INITIAL_REPLICA_COUNT is empty"
[ "$INITIAL_REPLICA_COUNT" = "$manifest_initial" ] || fail "INITIAL_REPLICA_COUNT ($INITIAL_REPLICA_COUNT) doesn't match the Deployment's original replicas ($manifest_initial)"

[ -n "${CONTAINER_IMAGE:-}" ] || fail "CONTAINER_IMAGE is empty"
[ "$CONTAINER_IMAGE" = "$live_image" ] || fail "CONTAINER_IMAGE ($CONTAINER_IMAGE) doesn't match the live container image ($live_image)"

[ -n "${FINAL_REPLICA_COUNT:-}" ] || fail "FINAL_REPLICA_COUNT is empty"
[ "$FINAL_REPLICA_COUNT" = "$REQUIRED_REPLICAS" ] || fail "FINAL_REPLICA_COUNT ($FINAL_REPLICA_COUNT) doesn't equal the required target ($REQUIRED_REPLICAS)"
[ "$FINAL_REPLICA_COUNT" = "$live_replicas" ] || fail "FINAL_REPLICA_COUNT doesn't match the Deployment's live replica count"

echo "PASS: deployment/hello is scaled to $REQUIRED_REPLICAS replicas and answers.env matches live cluster state"
