#!/usr/bin/env bash
set -euo pipefail
# Start nginx normally, then take away its ability to read the homepage
nginx -t
nginx
chmod 000 /var/www/html/index.html
