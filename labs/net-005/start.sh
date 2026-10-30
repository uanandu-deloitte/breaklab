#!/usr/bin/env bash
set -euo pipefail
# The backend API, listening locally on 8080. There is deliberately no
# /etc/hosts entry for api.internal yet — that's the broken state.
python3 /opt/app/app.py &
disown
exec sleep infinity
