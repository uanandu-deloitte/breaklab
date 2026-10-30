#!/usr/bin/env bash
set -uo pipefail
if [ -f /srv/app/VERSION ] && grep -q "v5.0.0" /srv/app/VERSION; then
  echo "PASS: v5.0.0 is deployed"
else
  echo "FAIL: v5.0.0 is not deployed"
  exit 1
fi
