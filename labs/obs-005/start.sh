#!/usr/bin/env bash
# PID 1. Seeds an initial app.log, starts a long-running reader that
# follows it by file descriptor (tail -f), then rotates the log out from
# under that reader exactly the way a real log-rotation job would: rename
# the old file aside, create a fresh empty one, and keep writing. The
# reader's descriptor stays pinned to the old (now renamed) file, so it
# never sees anything written to the new app.log.
set -u

mkdir -p /var/log
: > /var/log/app.log
for i in 1 2 3 4 5; do
  echo "request $i" >> /var/log/app.log
done

touch /var/log/reader-output.log
tail -f /var/log/app.log >> /var/log/reader-output.log 2>&1 &
reader_pid=$!
echo "$reader_pid" > /var/run/reader.pid

# Wait until the reader has genuinely opened its descriptor on app.log
# before rotating it out from underneath it — poll instead of a fixed
# sleep so this is correct regardless of how fast (or slow) the host
# happens to schedule the background process.
opened=0
for _ in $(seq 1 100); do
  for fd in /proc/"$reader_pid"/fd/*; do
    tgt=$(readlink -f "$fd" 2>/dev/null || true)
    if [ "$tgt" = "/var/log/app.log" ]; then
      opened=1
      break 2
    fi
  done
  sleep 0.05
done
[ "$opened" = "1" ] || echo "start.sh: warning: reader never appeared to open app.log" >&2

mv /var/log/app.log /var/log/app.log.1
touch /var/log/app.log
for i in 6 7 8 9 10; do
  echo "request $i" >> /var/log/app.log
done

exec sleep infinity
