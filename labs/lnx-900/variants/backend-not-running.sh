#!/usr/bin/env bash
set -euo pipefail
# Start everything normally, then kill the backend application process.
su -s /bin/bash -c "nohup python3 /opt/app/app.py >/tmp/app-backend.log 2>&1 &" appuser
sleep 1
nginx -t
nginx
sleep 1
pkill -f "python3 /opt/app/app.py" || true
