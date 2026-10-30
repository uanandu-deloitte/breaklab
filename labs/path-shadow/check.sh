#!/usr/bin/env bash
set -uo pipefail
deploy-sync
expected=$(/usr/bin/jq -r '.version' /etc/app/release.json)
actual=$(cat /var/lib/app/active-version.txt 2>/dev/null)
if [ "$actual" = "$expected" ]; then
  echo "PASS: active-version.txt matches release.json ($expected)"
else
  echo "FAIL: active-version.txt says '$actual', release.json says '$expected'"
  exit 1
fi
