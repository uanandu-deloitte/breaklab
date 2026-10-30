#!/usr/bin/env bash
set -euo pipefail
# Oops: a "simplify the paths" pass made the data directory relative
sed -i 's#DATA_DIR=/var/app/data#DATA_DIR=data#' /usr/local/bin/backup.sh
# Today's data changed; tonight's cron run was supposed to pick it up
echo v2 > /var/app/data/marker.txt
chown opsuser:opsuser /var/app/data/marker.txt
# Simulate cron actually invoking tonight's job: a real cron job always
# runs with the user's home directory as its working directory, not
# wherever the script happens to live
su - opsuser -c 'env -i HOME=/home/opsuser LOGNAME=opsuser PATH=/usr/bin:/usr/local/bin:/bin SHELL=/bin/sh sh -c "cd \"\$HOME\" && backup.sh"' >> /var/log/backup.log 2>&1 || true
