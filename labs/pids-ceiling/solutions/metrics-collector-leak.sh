#!/usr/bin/env bash
# Kills by scanning /proc with only bash builtins (no forking) -- pkill
# itself needs a fork to start, which the leak may not leave room for.
set -uo pipefail
for p in /proc/[0-9]*; do
  mapfile -d '' -t args < "$p/cmdline" 2>/dev/null || continue
  joined="${args[*]}"
  if [[ "$joined" == *metrics-collector.sh* ]]; then
    kill -9 "${p#/proc/}" 2>/dev/null
  fi
done
