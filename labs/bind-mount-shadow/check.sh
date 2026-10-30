#!/usr/bin/env bash
set -uo pipefail
appuser_uid=$(id -u appuser)
mount_uid=$(stat -c '%u' /data)
if [ "$mount_uid" -eq "$appuser_uid" ]; then
  echo "PASS: /data is owned by appuser (uid $appuser_uid)"
else
  echo "FAIL: /data is owned by uid $mount_uid, not appuser (uid $appuser_uid)"
  exit 1
fi
