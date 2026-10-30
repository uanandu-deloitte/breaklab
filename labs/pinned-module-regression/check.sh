#!/usr/bin/env bash
set -uo pipefail
expected=$(sed -n 's/^message: *//p' /infra/site.yaml)
actual=$(cat /etc/site/deployed.conf 2>/dev/null)
if [ "$expected" = "$actual" ]; then
  echo "PASS: /etc/site/deployed.conf matches /infra/site.yaml"
else
  echo "FAIL: deployed config ('${actual:-nothing}') doesn't match /infra/site.yaml ('$expected')"
  exit 1
fi
