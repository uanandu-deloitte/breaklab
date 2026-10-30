#!/usr/bin/env bash
set -euo pipefail
# Oops: the app's own logging config was pointed at a file path instead
# of stdout (someone copied a snippet from a non-containerized deploy).
mkdir -p /var/log/app
printf 'log_target=/var/log/app/app.log\n' > /etc/app/logging.conf
kill -HUP 1
sleep 1
