#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

TARGET_IMAGE="nginx:1.27-alpine"

status=$(kubectl rollout status deployment/web --timeout=10s 2>&1)
echo "$status" | grep -qi "successfully rolled out" || fail "rollout not complete: $status"

# kubectl rollout status can report success a moment before an old
# ReplicaSet's pod has actually finished terminating (pod deletion is
# async, with a grace period) — kubectl get pods can still briefly list
# it with the old image. Retry for a few seconds so that transient
# window doesn't fail a rollout that genuinely completed.
all_match=""
last_images=""
for i in $(seq 1 12); do
  images=$(kubectl get pods -l app=web --field-selector=status.phase=Running -o jsonpath='{range .items[*]}{.spec.containers[0].image}{"\n"}{end}' 2>/dev/null)
  last_images="$images"
  if [ -n "$images" ]; then
    all_match="yes"
    while IFS= read -r img; do
      [ -z "$img" ] && continue
      [ "$img" = "$TARGET_IMAGE" ] || all_match="no"
    done <<< "$images"
  else
    all_match="no"
  fi
  [ "$all_match" = "yes" ] && break
  sleep 1
done
[ "$all_match" = "yes" ] || fail "not every running pod for deployment/web runs $TARGET_IMAGE yet: $last_images"

ready=$(kubectl get deployment web -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
desired=$(kubectl get deployment web -o jsonpath='{.spec.replicas}' 2>/dev/null)
[ -n "$ready" ] && [ -n "$desired" ] && [ "$ready" = "$desired" ] || fail "deployment/web readyReplicas ($ready) does not match desired replicas ($desired)"

echo "PASS: every pod for deployment/web runs $TARGET_IMAGE and the rollout completed"
