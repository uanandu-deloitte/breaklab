#!/usr/bin/env bash
set -euo pipefail
python3 /opt/app/app.py &
disown
exec sleep infinity
