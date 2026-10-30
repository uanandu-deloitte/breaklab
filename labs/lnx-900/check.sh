#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8080}"
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") || {
  echo "FAIL: site unreachable or returned an HTTP error"; exit 1; }
if grep -q "BreakLab OK" <<<"$body"; then
  echo "PASS: customer portal serves the expected response"
else
  echo "FAIL: site is up but serving the wrong content"; exit 1
fi
