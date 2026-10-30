#!/usr/bin/env bash
set -uo pipefail
budget=$((50*1024*1024))
used=$(du -sb /var/app/data 2>/dev/null | awk '{print $1}')
if [ -z "$used" ]; then
  echo "FAIL: /var/app/data missing"
  exit 1
fi
if [ "$used" -lt "$budget" ]; then
  echo "PASS: /var/app/data is $((used/1024/1024))MB, under budget"
else
  echo "FAIL: /var/app/data is $((used/1024/1024))MB, over the 50MB budget"
  exit 1
fi
