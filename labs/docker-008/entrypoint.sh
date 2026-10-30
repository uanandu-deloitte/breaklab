#!/usr/bin/env bash
set -e
# The backend really does listen on 5001. The shipped nginx config's
# proxy_pass points at 5002 instead -- that mismatch is this lab's bug.
cd /srv/backend
nohup python3 -m http.server 5001 >/var/log/backend.log 2>&1 &
disown
nginx
exec sleep infinity
