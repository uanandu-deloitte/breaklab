#!/usr/bin/env bash
# Stand-in for a long-running service that keeps an append-only log file
# open for its entire lifetime (fd 3) and writes a heartbeat into it
# every couple seconds.
set -u
LOGFILE=/var/log/app/service.log
mkdir -p "$(dirname "$LOGFILE")"
exec 3>>"$LOGFILE"
echo $$ > /var/run/app.pid

while true; do
  echo "$(date -u +%FT%TZ) heartbeat" >&3
  # A plain `sleep 2` here would fork an external process that inherits
  # fd 3 (open file descriptors are copied across fork, not just exec),
  # so killing this daemon alone wouldn't release the deleted file's
  # space until that lingering child's own sleep also finished. A read
  # with a timeout is a bash builtin -- no fork, no extra fd holder.
  read -r -t 2 </dev/null || true
done
