#!/usr/bin/env bash
set -euo pipefail
nohup /usr/local/bin/metrics-collector.sh >/var/log/metrics-collector.log 2>&1 &
disown
