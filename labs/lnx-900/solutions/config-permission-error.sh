#!/usr/bin/env bash
set -euo pipefail
# Reference fix: least-privilege read access, then start the backend.
chgrp appuser /opt/app/config.yaml
chmod 640 /opt/app/config.yaml
su -s /bin/bash -c "nohup python3 /opt/app/app.py >/tmp/app-backend.log 2>&1 &" appuser
sleep 1
curl -fsS http://127.0.0.1:5000/ >/dev/null
