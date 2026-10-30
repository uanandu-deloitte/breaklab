#!/usr/bin/env bash
set -uo pipefail
version=$(cat /srv/app/current/VERSION 2>/dev/null)
if [ "$version" = "v2.0.0" ]; then
  echo "PASS: v2.0.0 is live"
else
  echo "FAIL: live version is ${version:-missing}, not v2.0.0"
  exit 1
fi
