#!/usr/bin/env bash
set -euo pipefail
# Reference fix: free the port, then start the real backend.
pkill -f "http.server 5000" || true
sleep 1
su -s /bin/bash -c "nohup python3 /opt/app/app.py >/tmp/app-backend.log 2>&1 &" appuser
sleep 1
curl -fsS http://127.0.0.1:5000/ >/dev/null
