#!/usr/bin/env bash
set -euo pipefail
echo "127.0.0.1 db.internal" >> /etc/hosts
nohup python3 /opt/db/db.py >/tmp/db.log 2>&1 &
sleep 1
nginx -t
nginx
# Make the api's own config unreadable to the account it runs as, then
# try to start it -- it will read-fail and exit immediately.
chown root:root /opt/api/config.yaml
chmod 600 /opt/api/config.yaml
su -s /bin/bash -c "nohup python3 /opt/api/api.py >/tmp/api.log 2>&1 &" apiuser || true
sleep 1
