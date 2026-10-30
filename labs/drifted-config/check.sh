#!/usr/bin/env bash
set -uo pipefail
message=$(sed -n 's/^message: *//p' /infra/site.yaml)
expected=$(MESSAGE="$message" envsubst '${MESSAGE}' < /infra/templates/index.html.tmpl)
actual=$(cat /var/www/html/index.html 2>/dev/null)
if [ "$expected" = "$actual" ]; then
  echo "PASS: deployed site matches /infra/site.yaml"
else
  echo "FAIL: deployed site has drifted from /infra/site.yaml"
  exit 1
fi
