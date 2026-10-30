#!/usr/bin/env bash
set -uo pipefail

pvc_phase=$(kubectl get pvc notes-pvc -o jsonpath='{.status.phase}' 2>/dev/null)
pod_phase=$(kubectl get pod notes-writer -o jsonpath='{.status.phase}' 2>/dev/null)

if [ "$pvc_phase" != "Bound" ]; then
  echo "FAIL: notes-pvc is ${pvc_phase:-missing}, not Bound"
  exit 1
fi
if [ "$pod_phase" != "Running" ]; then
  echo "FAIL: notes-writer is ${pod_phase:-missing}, not Running"
  exit 1
fi

marker="breaklab-check-$$"
if ! kubectl exec notes-writer -- sh -c "echo $marker > /notes/testfile" >/dev/null 2>&1; then
  echo "FAIL: could not write to /notes inside notes-writer"
  exit 1
fi
readback=$(kubectl exec notes-writer -- cat /notes/testfile 2>/dev/null)

if [ "$readback" = "$marker" ]; then
  echo "PASS: notes-pvc is Bound and genuinely mountable"
else
  echo "FAIL: wrote '$marker' to /notes/testfile but read back '${readback:-<empty>}'"
  exit 1
fi
