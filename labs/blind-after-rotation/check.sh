#!/usr/bin/env bash
set -uo pipefail
count=$(cat /var/lib/agent/metrics.count 2>/dev/null || echo "")
if [ "$count" = "8" ]; then
  echo "PASS: metrics agent reports the current log's 8 requests"
else
  echo "FAIL: metrics agent reports '$count', expected 8"
  exit 1
fi
