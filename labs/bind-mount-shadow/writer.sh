#!/usr/bin/env bash
set -uo pipefail
while true; do
  date +%s > /data/heartbeat.txt 2>>/var/log/writer-errors.log
  sleep 1
done
