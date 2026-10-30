#!/usr/bin/env bash
set -euo pipefail
# Oops: the agent's config was never updated after last night's log rotation
sed -i 's#LOGFILE=.*#LOGFILE=/var/log/app/app.log.1#' /etc/agent/agent.conf
collect-metrics
