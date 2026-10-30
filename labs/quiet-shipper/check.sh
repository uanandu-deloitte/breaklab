#!/usr/bin/env bash
set -uo pipefail
count=$(wc -l < /var/log/shipped.log 2>/dev/null || echo 0)
if [ "$count" -eq 10 ]; then
  echo "PASS: all 10 log lines were shipped"
else
  echo "FAIL: only $count/10 log lines reached the aggregator"
  exit 1
fi
