#!/usr/bin/env bash
set -euo pipefail
phase=$(kubectl get pvc reports-pvc -o jsonpath='{.status.phase}')
requested=$(kubectl get pvc reports-pvc -o jsonpath='{.spec.storageClassName}')
available=$(kubectl get storageclass -o jsonpath='{.items[*].metadata.name}' | tr ' ' '\n' | sort | tr '\n' ' ' | sed 's/ *$//')
pod_phase=$(kubectl get pod reports-writer -o jsonpath='{.status.phase}')

exists="no"
for sc in $available; do
  [ "$sc" = "$requested" ] && exists="yes"
done

mkdir -p /workspace
cat <<EOF > /workspace/answers.env
PVC_PHASE=$phase
REQUESTED_STORAGECLASS=$requested
AVAILABLE_STORAGECLASSES=$available
CLASS_EXISTS=$exists
POD_PHASE=$pod_phase
EOF
