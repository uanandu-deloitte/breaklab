#!/usr/bin/env bash
set -uo pipefail
render-dashboard
content=$(cat /var/lib/dashboards/overview.rendered 2>/dev/null || echo "no output")
if [ "$content" = "panel=requests value=42" ]; then
  echo "PASS: dashboard panel shows the current metric value"
else
  echo "FAIL: dashboard panel shows '$content', expected panel=requests value=42"
  exit 1
fi
