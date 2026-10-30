#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8080}"
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") || {
  echo "FAIL: site unreachable or returned an HTTP error"; exit 1; }
if grep -q "BreakLab net-002 OK" <<<"$body"; then
  echo "PASS: nginx is serving the real marker page on port 80"
else
  echo "FAIL: something is answering on port 80, but it isn't nginx serving the real page"
  exit 1
fi
