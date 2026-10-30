#!/usr/bin/env bash
set -uo pipefail
cd /srv/repo || { echo "FAIL: repo missing"; exit 1; }
actual=$(sed -n 's/^timeout_seconds: *//p' CONFIG.md 2>/dev/null)
if [ "$actual" = "60" ]; then
  echo "PASS: CONFIG.md has the hotfix (timeout_seconds: 60)"
else
  echo "FAIL: CONFIG.md still shows timeout_seconds: ${actual:-missing}, expected 60"
  exit 1
fi
