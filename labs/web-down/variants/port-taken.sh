#!/usr/bin/env bash
set -euo pipefail
# Grab port 80 first so nginx can't bind it
nohup python3 -m http.server 80 >/tmp/port-taken.log 2>&1 &
disown
sleep 1
nginx 2>/dev/null || true
