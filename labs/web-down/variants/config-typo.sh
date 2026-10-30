#!/usr/bin/env bash
set -euo pipefail
# Drop the semicolon so nginx refuses to start
sed -i 's/listen 80 default_server;/listen 80 default_server/' \
  /etc/nginx/sites-available/default
nginx 2>/dev/null || true
