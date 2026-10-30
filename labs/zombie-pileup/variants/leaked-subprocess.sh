#!/usr/bin/env bash
set -euo pipefail
# Oops: the worker pool launches jobs via subprocess.Popen and never
# waits on the ones that finish
nohup python3 /usr/local/bin/worker_pool.py >/var/log/worker_pool.log 2>&1 &
sleep 2
