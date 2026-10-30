#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point the proxy back at the port the backend actually listens on
sed -i 's/127.0.0.1:8082/127.0.0.1:8081/' /etc/nginx/sites-available/default
nginx -t
nginx -s reload
