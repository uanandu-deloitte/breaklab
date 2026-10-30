#!/usr/bin/env bash
set -euo pipefail
# Bring up the database and nginx normally.
nohup python3 /opt/db/db.py >/tmp/db.log 2>&1 &
sleep 1
nginx -t
nginx
# Start the api normally too -- but the hosts entry it needs to resolve
# db.internal was never added, so every request it handles will fail
# the moment it tries to reach the database.
su -s /bin/bash -c "nohup python3 /opt/api/api.py >/tmp/api.log 2>&1 &" apiuser
sleep 1
