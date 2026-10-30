#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the absolute path in the crontab entry
crontab -u opsuser -l | sed 's#[^ ]*backup\.sh#/usr/local/bin/backup.sh#' | crontab -u opsuser -
su - opsuser -c '/usr/local/bin/backup.sh'
