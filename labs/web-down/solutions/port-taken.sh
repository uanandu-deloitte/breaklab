#!/usr/bin/env bash
set -euo pipefail
# Reference fix: free the port and start nginx for real
pkill -f "http.server 80" || true
sleep 1
nginx -t
nginx
