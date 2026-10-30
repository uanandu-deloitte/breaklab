#!/usr/bin/env bash
set -uo pipefail
latest=$(cat /srv/artifacts/PROD_LATEST 2>/dev/null)
expected=$(tar -xOf "/srv/artifacts/$latest" CONTENT 2>/dev/null)
actual=$(cat /srv/app/live/CONTENT 2>/dev/null)
if [ -n "$expected" ] && [ "$expected" = "$actual" ]; then
  echo "PASS: live artifact matches the current prod build ($latest)"
else
  echo "FAIL: live artifact ('${actual:-missing}') doesn't match current prod build $latest ('${expected:-missing}')"
  exit 1
fi
