#!/usr/bin/env bash
# Writes the same kind of event to two places at once: stdout (which
# docker logs captures) and a file on disk (which it never will).
i=0
while true; do
  i=$((i + 1))
  echo "heartbeat $i"
  echo "file-only $i $(date -Is)" >> /var/log/app/app.log
  sleep 1
done
