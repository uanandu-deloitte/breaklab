#!/usr/bin/env bash
set -euo pipefail
# Start the backend normally, then point nginx's upstream at the wrong port.
su -s /bin/bash -c "nohup python3 /opt/app/app.py >/tmp/app-backend.log 2>&1 &" appuser
sleep 1
sed -i 's/proxy_pass http:\/\/127.0.0.1:5000;/proxy_pass http:\/\/127.0.0.1:5999;/' \
  /etc/nginx/sites-available/default
nginx -t
nginx
