#!/usr/bin/env bash
set -e
# A decoy process grabs port 80 first, so nginx fails to bind it.
nohup python3 -m http.server 80 >/tmp/decoy.log 2>&1 &
disown
sleep 1
nginx 2>/dev/null || true
exec sleep infinity
