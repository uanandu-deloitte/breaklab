#!/usr/bin/env bash
set -euo pipefail
# Oops: a "log cleanup" cron job rm's the live log directly instead of
# rotating it (no copytruncate, no signal to the app to reopen).
for i in $(seq 1 10); do
  [ -s /var/run/app.pid ] && kill -0 "$(cat /var/run/app.pid)" 2>/dev/null && break
  sleep 0.3
done
rm -f /var/log/app/service.log
