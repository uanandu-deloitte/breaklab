#!/usr/bin/env bash
set -euo pipefail
echo "127.0.0.1 db.internal" >> /etc/hosts
# Deliberately never start the database -- the api and nginx are both
# otherwise healthy.
nginx -t
nginx
su -s /bin/bash -c "nohup python3 /opt/api/api.py >/tmp/api.log 2>&1 &" apiuser
sleep 1
