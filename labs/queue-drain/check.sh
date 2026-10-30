#!/usr/bin/env bash
set -uo pipefail
count=$(ls -1 /var/spool/jobs/done 2>/dev/null | wc -l)
if [ "$count" -eq 5 ]; then
  echo "PASS: all 5 jobs drained"
else
  echo "FAIL: only $count/5 jobs drained"
  exit 1
fi
