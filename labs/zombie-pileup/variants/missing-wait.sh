#!/usr/bin/env bash
set -euo pipefail
# Oops: the job supervisor never reaps the children it forks
nohup python3 /usr/local/bin/supervisor.py >/var/log/supervisor.log 2>&1 &
sleep 2
