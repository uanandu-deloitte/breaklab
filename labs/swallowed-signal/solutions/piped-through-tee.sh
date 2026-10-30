#!/usr/bin/env bash
set -euo pipefail
# Reference fix: drop the pipe so the app process itself can take PID
# 1's direct child slot and receive signals directly.
printf 'exec python3 /usr/local/bin/app.py\n' > /etc/app/launch-cmd
kill -HUP 1
sleep 0.5
