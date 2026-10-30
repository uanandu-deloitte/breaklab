#!/usr/bin/env bash
set -euo pipefail
# Oops: someone repointed clocksync at a renamed internal host that no
# longer exists, so syncs "succeed" but never actually correct anything
printf '%s' 'ntp.old-internal.example' > /etc/clocksync/server
echo -14400 > /var/lib/clock/offset_seconds
log-event "still running, sync reports healthy"
