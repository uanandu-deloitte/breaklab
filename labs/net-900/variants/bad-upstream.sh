#!/usr/bin/env bash
set -euo pipefail
echo "127.0.0.1 db.internal" >> /etc/hosts
nohup python3 /opt/db/db.py >/tmp/db.log 2>&1 &
sleep 1
su -s /bin/bash -c "nohup python3 /opt/api/api.py >/tmp/api.log 2>&1 &" apiuser
sleep 1
# Point nginx's upstream at the wrong port for the api process.
sed -i 's/proxy_pass http:\/\/127.0.0.1:5000;/proxy_pass http:\/\/127.0.0.1:5999;/' \
  /etc/nginx/sites-available/default
nginx -t
nginx
