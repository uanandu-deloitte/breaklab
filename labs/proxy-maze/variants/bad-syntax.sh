#!/usr/bin/env bash
set -euo pipefail
cd /srv/backend && nohup python3 -m http.server 8081 >/tmp/backend.log 2>&1 &
sleep 1
# Oops: a manual edit dropped a closing brace
sed -i '/# location-close/d' /etc/nginx/sites-available/default
nginx 2>/dev/null || true
