#!/usr/bin/env bash
# Runs on the host: confirms `docker logs` is actually receiving live
# output from the running app, not just that the app is alive.
set -uo pipefail
container="${BREAKLAB_CONTAINER:?}"

before=$(docker logs "$container" 2>&1 | grep -c '^heartbeat ' || true)
sleep 3
after=$(docker logs "$container" 2>&1 | grep -c '^heartbeat ' || true)

if [ "$after" -gt "$before" ]; then
  echo "PASS: docker logs is receiving live output from the app ($before -> $after)"
else
  echo "FAIL: docker logs shows no new output from the running app ($before -> $after)"
  exit 1
fi
