#!/usr/bin/env bash
set -uo pipefail
ip=$(kubectl get svc orders -o jsonpath='{.spec.clusterIP}' 2>/dev/null)
port=$(kubectl get svc orders -o jsonpath='{.spec.ports[0].port}' 2>/dev/null)
if [ -z "$ip" ] || [ -z "$port" ]; then
  echo "FAIL: orders Service not found"
  exit 1
fi
if curl -fsS --max-time 3 "http://${ip}:${port}/" >/dev/null 2>&1; then
  echo "PASS: orders Service routes to a healthy backend"
else
  echo "FAIL: orders Service is not routing traffic to any backend"
  exit 1
fi
