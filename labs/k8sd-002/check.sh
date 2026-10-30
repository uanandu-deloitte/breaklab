#!/usr/bin/env bash
set -uo pipefail

qos=$(kubectl get pod qos-demo -o jsonpath='{.status.qosClass}' 2>/dev/null)
phase=$(kubectl get pod qos-demo -o jsonpath='{.status.phase}' 2>/dev/null)

if [ "$qos" = "Guaranteed" ] && [ "$phase" = "Running" ]; then
  echo "PASS: qos-demo is Running with Guaranteed QoS"
else
  echo "FAIL: qos-demo is phase=${phase:-missing}, qosClass=${qos:-missing} (need Running + Guaranteed)"
  exit 1
fi
