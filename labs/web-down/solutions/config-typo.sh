#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the semicolon and start nginx
sed -i 's/listen 80 default_server$/listen 80 default_server;/' \
  /etc/nginx/sites-available/default
nginx -t
nginx
