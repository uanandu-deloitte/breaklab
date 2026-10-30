#!/usr/bin/env bash
set -uo pipefail
if [ -f /srv/app/VERSION ] && grep -q "v2.0.0" /srv/app/VERSION; then
  echo "PASS: v2.0.0 is deployed"
else
  echo "FAIL: /srv/app is not on v2.0.0"
  exit 1
fi
