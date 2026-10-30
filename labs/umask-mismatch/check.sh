#!/usr/bin/env bash
set -uo pipefail
rm -f /var/lib/appdata/export.json
su -s /bin/bash producer -c /usr/local/bin/run-producer.sh

if [ ! -f /var/lib/appdata/export.json ]; then
  echo "FAIL: producer did not write export.json"
  exit 1
fi

if su -s /bin/sh consumer -c 'cat /var/lib/appdata/export.json' >/dev/null 2>&1; then
  echo "PASS: consumer can read the latest export"
else
  echo "FAIL: consumer cannot read /var/lib/appdata/export.json"
  exit 1
fi
