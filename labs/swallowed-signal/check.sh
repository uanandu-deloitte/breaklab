#!/usr/bin/env bash
# Runs on the host: sends SIGTERM straight to the container's PID 1 the
# same way `docker stop` eventually would, then checks (fast, without
# waiting for a real stop timeout) whether the app actually saw it.
set -uo pipefail
container="${BREAKLAB_CONTAINER:?}"
marker=/var/run/app/graceful-shutdown.marker

docker exec "$container" rm -f "$marker" >/dev/null 2>&1 || true
docker kill -s TERM "$container" >/dev/null 2>&1

ok=0
for i in $(seq 1 10); do
  if docker exec "$container" test -f "$marker" >/dev/null 2>&1; then
    ok=1
    break
  fi
  sleep 0.5
done

if [ "$ok" -eq 1 ]; then
  echo "PASS: app received SIGTERM and shut down gracefully"
else
  echo "FAIL: app never saw SIGTERM — graceful shutdown marker not written"
  exit 1
fi
