#!/usr/bin/env bash
# Edits and kills use only bash builtins (no forking) so the fix still
# applies even while the box is still short on pids; only the final
# restart needs one real fork, by which point the kills above have
# freed plenty of room. Stopping flood.sh matters too: it's standing in
# for a client retry storm caused by the stuck jobs, which would ease
# off once the handler starts responding again.
set -uo pipefail
content=$(< /usr/local/bin/handler.py)
printf '%s\n' "${content//CAP = 999999/CAP = 10}" > /usr/local/bin/handler.py

for p in /proc/[0-9]*; do
  mapfile -d '' -t args < "$p/cmdline" 2>/dev/null || continue
  joined="${args[*]}"
  if [[ "$joined" == *handler.py* || "$joined" == *flood.sh* || "$joined" == *"jobs/done"* ]]; then
    kill -9 "${p#/proc/}" 2>/dev/null
  fi
done

rm -f /var/spool/jobs/*.job

sleep 0.5
nohup python3 /usr/local/bin/handler.py >/var/log/handler.log 2>&1 &
disown
