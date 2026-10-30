#!/usr/bin/env bash
set -uo pipefail
sync-clock >/dev/null
log-event "checking in" >/dev/null
line=$(tail -1 /var/log/host-a.log)
ts=$(awk '{print $2}' <<<"$line")
now=$(date +%s)
if [ "$ts" -ge "$now" ]; then
  diff=$((ts - now))
else
  diff=$((now - ts))
fi
if [ "$diff" -le 5 ]; then
  echo "PASS: host-a's clock is in sync (within ${diff}s of real time)"
else
  echo "FAIL: host-a's clock is off by ${diff}s"
  exit 1
fi
