#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore a valid config and recollect
printf 'LOGFILE=/var/log/app/app.log\n' > /etc/agent/agent.conf
collect-metrics
