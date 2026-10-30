#!/usr/bin/env bash
set -euo pipefail
# Oops: the job that's supposed to drain the queue on a timer hasn't run
# in a while, so the backlog built up and filled the queue entirely
for i in $(seq 1 5); do
  printf 'backlog item %d\n' "$i" > "/var/spool/shipper/queue/msg-backlog-$i"
done
