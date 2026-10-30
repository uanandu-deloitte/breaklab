#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

POD=scratch-writer

phase=$(kubectl get pod "$POD" -o jsonpath='{.status.phase}' 2>/dev/null)
[ "$phase" = "Running" ] || fail "pod $POD is not Running (phase: $phase)"

ready=$(kubectl get pod "$POD" -o jsonpath='{.status.containerStatuses[0].ready}' 2>/dev/null)
[ "$ready" = "true" ] || fail "pod $POD's container is not Ready"

runAsNonRoot=$(kubectl get pod "$POD" -o jsonpath='{.spec.securityContext.runAsNonRoot}' 2>/dev/null)
[ "$runAsNonRoot" = "true" ] || fail "spec.securityContext.runAsNonRoot is not set to true (got: $runAsNonRoot)"

runAsUser=$(kubectl get pod "$POD" -o jsonpath='{.spec.securityContext.runAsUser}' 2>/dev/null)
[ "$runAsUser" = "10001" ] || fail "spec.securityContext.runAsUser is not 10001 (got: $runAsUser)"

readOnly=$(kubectl get pod "$POD" -o jsonpath='{.spec.containers[0].securityContext.readOnlyRootFilesystem}' 2>/dev/null)
[ "$readOnly" = "true" ] || fail "containers[0].securityContext.readOnlyRootFilesystem is not set to true (got: $readOnly)"

live_uid=$(kubectl exec "$POD" -- id -u 2>/dev/null)
[ "$live_uid" = "10001" ] || fail "the container is actually running as uid $live_uid, not 10001"

# It must still be able to do its job: write to its designated scratch
# path, backed by the emptyDir volume, and read it back.
marker="check-$$"
kubectl exec "$POD" -- sh -c "echo $marker > /data/checkfile" >/dev/null 2>&1 \
  || fail "could not write to /data/checkfile — the scratch volume isn't mounted or isn't writable"
readback=$(kubectl exec "$POD" -- cat /data/checkfile 2>/dev/null)
[ "$readback" = "$marker" ] || fail "/data/checkfile did not round-trip the written marker (got: $readback)"

# The true root filesystem (outside the emptyDir) must actually be
# read-only — proves this wasn't satisfied by just adding the
# securityContext without also adding a writable volume for /data.
if kubectl exec "$POD" -- sh -c "touch /root-write-test" >/dev/null 2>&1; then
  fail "the container's root filesystem is still writable — readOnlyRootFilesystem isn't really in effect"
fi

echo "PASS: scratch-writer runs non-root, read-only rootfs, and still writes to its scratch volume"
