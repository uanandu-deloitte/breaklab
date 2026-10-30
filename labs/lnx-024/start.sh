#!/usr/bin/env bash
set -euo pipefail
/opt/app/worker.sh &
exec python3 -m http.server 80 --directory /srv/www
