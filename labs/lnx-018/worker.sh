#!/usr/bin/env bash
# A tiny worker that traps SIGTERM and writes a completion marker before
# exiting — but only on SIGTERM. SIGKILL gives it no chance to run this
# trap at all, which is the entire point of this lab.
set -u
marker="${1:?usage: worker.sh <marker-path>}"

trap 'echo "graceful shutdown at $(date -Iseconds)" > "$marker"; exit 0' TERM

while true; do
  sleep 1
done
