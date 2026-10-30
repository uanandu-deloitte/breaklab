#!/usr/bin/env bash
set -euo pipefail
# Oops: clock sync got disabled during a maintenance window and nobody
# turned it back on; the clock has been free-running since
printf '%s' 'false' > /etc/clocksync/enabled
echo -14400 > /var/lib/clock/offset_seconds
log-event "still running, clock drifting"
