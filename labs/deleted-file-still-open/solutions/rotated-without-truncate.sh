#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restart the daemon so it reopens a fresh log file,
# dropping its handle on the deleted one.
pid=$(cat /var/run/app.pid 2>/dev/null || true)
if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
  kill "$pid"
  for i in $(seq 1 20); do
    kill -0 "$pid" 2>/dev/null || break
    sleep 0.2
  done
fi
nohup /usr/local/bin/logger-daemon.sh >/var/log/app/daemon.out 2>&1 &
disown
sleep 1
