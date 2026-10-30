#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8081}"
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") || {
  echo "FAIL: storefront proxy unreachable or returned an HTTP error"; exit 1; }
if grep -q "Storefront v2" <<<"$body"; then
  echo "PASS: storefront proxy is serving the backend correctly"
else
  echo "FAIL: storefront proxy is up but serving the wrong content"; exit 1
fi
