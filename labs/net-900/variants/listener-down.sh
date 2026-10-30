#!/usr/bin/env bash
set -euo pipefail
echo "127.0.0.1 db.internal" >> /etc/hosts
nohup python3 /opt/db/db.py >/tmp/db.log 2>&1 &
sleep 1
nginx -t
nginx
su -s /bin/bash -c "nohup python3 /opt/api/api.py >/tmp/api.log 2>&1 &" apiuser
sleep 1
# Confirmed running -- now kill it, same as a crashed/stopped service.
pkill -f "python3 /opt/api/api.py" || true
sleep 1
