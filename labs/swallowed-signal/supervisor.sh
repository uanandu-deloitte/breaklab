#!/usr/bin/env bash
# PID 1. Spawns the app process according to /etc/app/launch-cmd,
# forwards SIGTERM it receives to that direct child, and can respawn
# the child on SIGHUP (used to re-apply a changed launch-cmd without
# restarting the whole container).
set -u
CONF=/etc/app/launch-cmd
CHILD_PID=""

spawn() {
  CHILD_PID=""
  cmd=$(cat "$CONF" 2>/dev/null)
  [ -z "$cmd" ] && return
  bash -c "$cmd" &
  CHILD_PID=$!
}

forward_term() {
  [ -n "$CHILD_PID" ] && kill -TERM "$CHILD_PID" 2>/dev/null
}

respawn() {
  if [ -n "$CHILD_PID" ]; then
    kill -KILL "$CHILD_PID" 2>/dev/null
  fi
  spawn
}

trap forward_term TERM
trap respawn HUP

spawn

while true; do
  wait -n 2>/dev/null || sleep 0.2
done
