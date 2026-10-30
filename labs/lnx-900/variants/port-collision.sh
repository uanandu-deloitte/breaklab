#!/usr/bin/env bash
set -euo pipefail
nginx -t
nginx
# Squat on the backend's port before it gets a chance to bind it.
nohup python3 -m http.server 5000 --bind 127.0.0.1 >/tmp/port-collision.log 2>&1 &
disown
sleep 1
su -s /bin/bash -c "nohup python3 /opt/app/app.py >/tmp/app-backend.log 2>&1 &" appuser || true
sleep 1
