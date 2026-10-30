#!/usr/bin/env bash
set -uo pipefail
expected=$(sed -n 's/^message: *//p' /infra/resources.yaml)
actual=$(cat /etc/motd.conf 2>/dev/null)
if [ "$expected" = "$actual" ]; then
  echo "PASS: /etc/motd.conf matches /infra/resources.yaml"
else
  echo "FAIL: infra-apply is skipping /etc/motd.conf instead of applying it"
  exit 1
fi
