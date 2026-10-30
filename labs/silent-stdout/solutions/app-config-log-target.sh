#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point the app's own logging config back at stdout.
printf 'log_target=stdout\n' > /etc/app/logging.conf
kill -HUP 1
sleep 1
