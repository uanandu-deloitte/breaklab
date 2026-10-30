#!/usr/bin/env bash
set -euo pipefail
# Oops: a crontab "cleanup" dropped the absolute path
crontab -u opsuser -l | sed 's#/usr/local/bin/backup.sh#backup.sh#' | crontab -u opsuser -
# Today's data changed; tonight's cron run was supposed to pick it up
echo v2 > /var/app/data/marker.txt
chown opsuser:opsuser /var/app/data/marker.txt
# Simulate cron actually invoking tonight's job, under cron's own minimal
# environment -- not an interactive login shell's
su - opsuser -c 'env -i HOME=/home/opsuser LOGNAME=opsuser PATH=/usr/bin:/bin SHELL=/bin/sh sh -c "backup.sh"' >> /var/log/backup.log 2>&1 || true
