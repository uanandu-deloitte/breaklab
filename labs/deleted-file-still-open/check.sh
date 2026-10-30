#!/usr/bin/env bash
set -uo pipefail

pidfile=/var/run/app.pid
pid=""
for i in $(seq 1 10); do
  if [ -s "$pidfile" ]; then
    pid=$(cat "$pidfile")
    kill -0 "$pid" 2>/dev/null && break
  fi
  pid=""
  sleep 0.3
done

if [ -z "$pid" ]; then
  echo "FAIL: logger-daemon is not running"
  exit 1
fi

leaked=0
for fd in /proc/[0-9]*/fd/*; do
  [ -e "$fd" ] || continue
  target=$(readlink "$fd" 2>/dev/null) || continue
  case "$target" in
    *"(deleted)"*) leaked=1 ;;
  esac
done

if [ "$leaked" -eq 1 ]; then
  echo "FAIL: a process is still holding a deleted file open — disk space not reclaimed"
  exit 1
fi

if [ ! -f /var/log/app/service.log ]; then
  echo "FAIL: service.log is missing — the app isn't logging"
  exit 1
fi

echo "PASS: no deleted file descriptors leaking space, app logging normally"
