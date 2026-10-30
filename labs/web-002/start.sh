#!/usr/bin/env bash
set -euo pipefail
# Backend: a tiny local web server with its own distinct page, on a port
# nginx does not (yet) proxy to.
cd /srv/backend && nohup python3 -m http.server 8090 >/tmp/backend.log 2>&1 &
# nginx: installed with its stock default site, untouched -- it only
# serves its own default page on :80 until its config is edited.
nginx
exec sleep infinity
