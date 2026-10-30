#!/usr/bin/env bash
set -euo pipefail
# Oops: an engineer's `tail -f` from an earlier debugging session was
# never closed, and the log got cleaned up out from under it too.
for i in $(seq 1 10); do
  [ -s /var/run/app.pid ] && kill -0 "$(cat /var/run/app.pid)" 2>/dev/null && break
  sleep 0.3
done
nohup tail -f /var/log/app/service.log >/dev/null 2>&1 &
echo $! > /var/run/stray.pid
disown
sleep 0.3
rm -f /var/log/app/service.log
