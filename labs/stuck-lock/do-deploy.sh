#!/usr/bin/env bash
set -euo pipefail
lockdir=/var/run/deploy
lock="$lockdir/deploy.lock"
if [ -e "$lock" ]; then
  echo "deploy: another deploy is already in progress ($lock exists)" >&2
  exit 1
fi
touch "$lock"
echo "v5.0.0" > /srv/app/VERSION
rm -f "$lock"
echo "deploy: v5.0.0 is live"
