#!/usr/bin/env bash
set -euo pipefail
sed -i 's/^CAP = .*/CAP = 999999/' /usr/local/bin/handler.py
pkill -f handler.py || true
sleep 0.5
nohup python3 /usr/local/bin/handler.py >/var/log/handler.log 2>&1 &
disown
nohup /usr/local/bin/flood.sh >/var/log/flood.log 2>&1 &
disown
