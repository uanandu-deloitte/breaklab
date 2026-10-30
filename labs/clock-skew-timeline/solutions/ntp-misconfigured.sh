#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point clocksync back at the real internal server
printf '%s' 'ntp.internal.example' > /etc/clocksync/server
sync-clock
