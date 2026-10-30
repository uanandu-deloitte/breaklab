#!/usr/bin/env bash
set -euo pipefail
# Oops: someone added "proper" file-based logging to the launch command
# itself, redirecting the app's own stdout/stderr into a file.
printf 'exec python3 /usr/local/bin/app.py >>/var/log/app/app.log 2>&1\n' > /etc/app/launch-cmd
kill -HUP 1
sleep 1
