#!/usr/bin/env bash
set -euo pipefail
# The real backend listens on 8095 -- but the nginx config shipped here
# points proxy_pass at 8098, a port nothing is bound to.
cd /srv/backend && nohup python3 -m http.server 8095 >/tmp/backend.log 2>&1 &
nginx
exec sleep infinity
