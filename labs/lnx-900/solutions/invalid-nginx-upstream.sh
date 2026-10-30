#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point nginx back at the backend's real address and reload.
sed -i 's/proxy_pass http:\/\/127.0.0.1:5999;/proxy_pass http:\/\/127.0.0.1:5000;/' \
  /etc/nginx/sites-available/default
nginx -t
nginx -s reload
