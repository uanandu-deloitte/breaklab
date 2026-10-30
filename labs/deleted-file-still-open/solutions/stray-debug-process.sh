#!/usr/bin/env bash
set -euo pipefail
# Reference fix: kill the leftover tail -f (not the app's fault at all),
# then restart the app itself too since its own handle is now stale.
if [ -f /var/run/stray.pid ]; then
  spid=$(cat /var/run/stray.pid)
  if kill -0 "$spid" 2>/dev/null; then
    kill "$spid"
  fi
  rm -f /var/run/stray.pid
fi

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
