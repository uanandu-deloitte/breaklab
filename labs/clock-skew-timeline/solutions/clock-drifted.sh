#!/usr/bin/env bash
set -euo pipefail
# Reference fix: turn sync back on and let it correct the clock
printf '%s' 'true' > /etc/clocksync/enabled
sync-clock
