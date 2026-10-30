#!/usr/bin/env bash
set -euo pipefail
# Reference fix: least-privilege read access, then start the api.
chgrp apiuser /opt/api/config.yaml
chmod 640 /opt/api/config.yaml
su -s /bin/bash -c "nohup python3 /opt/api/api.py >/tmp/api.log 2>&1 &" apiuser
sleep 1
curl -fsS http://127.0.0.1:5000/ >/dev/null
