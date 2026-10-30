#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8080}"
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") || {
  echo "FAIL: site unreachable through the mapped host port ${port}"; exit 1; }
if grep -q "BreakLab docker-006 OK" <<<"$body"; then
  echo "PASS: the mapped host port reaches the server"
else
  echo "FAIL: something answered on ${port}, but not with the expected page"
  exit 1
fi
