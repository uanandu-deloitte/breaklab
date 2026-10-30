#!/usr/bin/env bash
set -euo pipefail
# Reference fix: drop the redirection so the app's output goes to the
# container's actual stdout again.
printf 'exec python3 /usr/local/bin/app.py\n' > /etc/app/launch-cmd
kill -HUP 1
sleep 1
