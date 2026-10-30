#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

SCRIPT=/workspace/health.sh
[ -f "$SCRIPT" ] || fail "/workspace/health.sh not found"
[ -x "$SCRIPT" ] || fail "/workspace/health.sh is not executable"

restore_worker() {
  pgrep -f worker.sh >/dev/null 2>&1 && return
  nohup /opt/app/worker.sh >/dev/null 2>&1 &
  disown 2>/dev/null || true
}

# Healthy fixture: health.sh should pass cleanly.
if ! "$SCRIPT" >/tmp/health-healthy.out 2>&1; then
  restore_worker
  fail "health.sh exited non-zero against the healthy fixture (see /tmp/health-healthy.out)"
fi

# Break the fixture: stop the background worker it's supposed to watch
# for, then confirm health.sh notices.
pkill -f worker.sh || true
sleep 1

if "$SCRIPT" >/tmp/health-unhealthy.out 2>&1; then
  restore_worker
  fail "health.sh exited 0 even though the worker process is gone"
fi

restore_worker
sleep 1

echo "PASS: health.sh correctly reports both the healthy and the unhealthy state"
