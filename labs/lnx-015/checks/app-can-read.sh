#!/usr/bin/env bash
set -uo pipefail
code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 3 "http://localhost:80/healthz")
if [ "$code" = "200" ]; then
  echo "PASS: health endpoint responds 200"
else
  echo "FAIL: health endpoint returned $code"
  exit 1
fi
