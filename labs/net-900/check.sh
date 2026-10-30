#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8080}"
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") || {
  echo "FAIL: public API unreachable or returned an HTTP error"; exit 1; }
if grep -q "net900-api-db-chain-healthy" <<<"$body"; then
  echo "PASS: public API is serving through a healthy nginx -> api -> database chain"
else
  echo "FAIL: public API responded, but the full chain isn't healthy"; exit 1
fi
