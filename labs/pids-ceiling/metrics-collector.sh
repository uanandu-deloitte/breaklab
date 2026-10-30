#!/usr/bin/env bash
# Takes a cheap process snapshot every 0.05s -- except the snapshot got
# wrapped in a leftover debugging delay that never got removed, so each
# snapshot lingers for 5s instead of finishing immediately. They pile up
# far faster than they drain.
while true; do
  ( sleep 5; ps aux >/var/log/metrics-snapshot.log 2>&1 ) &
  sleep 0.05
done
