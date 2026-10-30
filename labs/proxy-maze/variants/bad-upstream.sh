#!/usr/bin/env bash
set -euo pipefail
cd /srv/backend && nohup python3 -m http.server 8081 >/tmp/backend.log 2>&1 &
sleep 1
# Oops: a "port renumbering" pass pointed the proxy at the wrong backend port
sed -i 's/127.0.0.1:8081/127.0.0.1:8082/' /etc/nginx/sites-available/default
nginx
