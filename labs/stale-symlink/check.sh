#!/usr/bin/env bash
set -uo pipefail
content=$(cat /opt/current/app.txt 2>/dev/null)
if [ "$content" = "v2" ]; then
  echo "PASS: /opt/current serves v2"
else
  echo "FAIL: /opt/current does not serve v2"
  exit 1
fi
