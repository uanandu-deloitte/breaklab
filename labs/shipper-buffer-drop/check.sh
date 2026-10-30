#!/usr/bin/env bash
set -uo pipefail
marker="probe-$$-$(date +%s)"
echo "$marker" >> /var/log/app.log
ship-logs
drain-queue
if grep -qF "$marker" /var/log/shipped-final.log 2>/dev/null; then
  echo "PASS: a fresh log line made it all the way to the aggregator"
else
  echo "FAIL: the fresh probe line never reached the aggregator"
  exit 1
fi
