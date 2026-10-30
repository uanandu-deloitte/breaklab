#!/usr/bin/env bash
set -euo pipefail
# Oops: a manual edit left the config file syntactically broken
rm -f /var/lib/agent/metrics.count
printf 'LOGFILE=/var/log/app/app.log\nBAD LINE WITHOUT EQUALS\n' > /etc/agent/agent.conf
