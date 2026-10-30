#!/usr/bin/env bash
set -euo pipefail
# Reference fix: the log doesn't need to be huge; rotate it back down
: > /var/app/data/app.log
echo "app data" >> /var/app/data/app.log
