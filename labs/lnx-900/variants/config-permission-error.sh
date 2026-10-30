#!/usr/bin/env bash
set -euo pipefail
nginx -t
nginx
# Make the backend's own config unreadable to the account it runs as,
# then try to start it — it will read-fail and exit immediately.
chown root:root /opt/app/config.yaml
chmod 600 /opt/app/config.yaml
su -s /bin/bash -c "nohup python3 /opt/app/app.py >/tmp/app-backend.log 2>&1 &" appuser || true
sleep 1
