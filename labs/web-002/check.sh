#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8080}"
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") || {
  echo "FAIL: proxy unreachable or returned an HTTP error"; exit 1; }
if grep -q "web002-backend-ok" <<<"$body"; then
  echo "PASS: nginx is proxying to the backend"
else
  echo "FAIL: nginx responded, but not with the backend's content"; exit 1
fi
