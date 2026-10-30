#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point the agent at today's active log and recollect
sed -i 's#LOGFILE=.*#LOGFILE=/var/log/app/app.log#' /etc/agent/agent.conf
collect-metrics
