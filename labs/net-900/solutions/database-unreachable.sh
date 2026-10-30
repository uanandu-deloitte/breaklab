#!/usr/bin/env bash
set -euo pipefail
# Reference fix: start the database.
nohup python3 /opt/db/db.py >/tmp/db.log 2>&1 &
sleep 1
curl -fsS http://127.0.0.1:5000/ >/dev/null
