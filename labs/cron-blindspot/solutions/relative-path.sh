#!/usr/bin/env bash
set -euo pipefail
# Reference fix: use an absolute path so cron's working directory can't matter
sed -i 's#DATA_DIR=data#DATA_DIR=/var/app/data#' /usr/local/bin/backup.sh
su - opsuser -c '/usr/local/bin/backup.sh'
