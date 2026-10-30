#!/usr/bin/env bash
set -uo pipefail
expected=$(cat /srv/app/requirements.lock 2>/dev/null)
actual=$(cat /srv/app/vendor/requirements.lock 2>/dev/null)
if [ -n "$expected" ] && [ "$expected" = "$actual" ]; then
  echo "PASS: deployed vendor matches requirements.lock ($expected)"
else
  echo "FAIL: deployed vendor ('${actual:-missing}') doesn't match requirements.lock ('${expected:-missing}')"
  exit 1
fi
