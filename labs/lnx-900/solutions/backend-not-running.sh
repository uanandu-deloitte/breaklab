#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restart the backend application process.
su -s /bin/bash -c "nohup python3 /opt/app/app.py >/tmp/app-backend.log 2>&1 &" appuser
sleep 1
curl -fsS http://127.0.0.1:5000/ >/dev/null
