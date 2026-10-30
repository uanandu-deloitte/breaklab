#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

for d in bin config data logs; do
  [ -d "/opt/acme-app/$d" ] || fail "/opt/acme-app/$d is missing"
done

[ -f /opt/acme-app/config/app.conf ] || fail "missing placeholder config file at /opt/acme-app/config/app.conf"
[ -f /opt/acme-app/logs/app.log ] || fail "missing placeholder log file at /opt/acme-app/logs/app.log"

[ -f /opt/fixtures/acme-app ] || fail "the original fixture at /opt/fixtures/acme-app is gone — copy it, don't move it"
[ -x /opt/acme-app/bin/acme-app ] || fail "/opt/acme-app/bin/acme-app is missing or not executable"
cmp -s /opt/fixtures/acme-app /opt/acme-app/bin/acme-app || fail "/opt/acme-app/bin/acme-app doesn't match the supplied fixture"

echo "PASS: /opt/acme-app layout matches the required structure"
