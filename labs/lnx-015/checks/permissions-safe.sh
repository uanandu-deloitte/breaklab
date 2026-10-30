#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

[ -f /srv/app/config.yaml ] || fail "/srv/app/config.yaml is missing"

want_content=$'service: acme\nenv: production\n'
got_content=$(cat /srv/app/config.yaml)
if [ "$got_content" != "${want_content%$'\n'}" ]; then
  fail "/srv/app/config.yaml content has changed — fix permissions, not the file's content"
fi

mode=$(stat -c '%a' /srv/app/config.yaml)
last_digit=${mode: -1}
if [ "$mode" = "777" ] || (( last_digit & 2 )); then
  fail "/srv/app/config.yaml is world-writable (mode $mode) — this is not allowed"
fi

# pgrep -x matches the process's exact name only (not its full argv) —
# -f would also match the "su -s /bin/bash -c python3 /opt/app.py
# appuser" wrapper process that starts the app as root before it drops
# to appuser, since that wrapper's own command line happens to contain
# the same substring.
app_pid=$(pgrep -x python3 | head -1)
if [ -z "$app_pid" ]; then
  fail "the application process isn't running"
fi
app_user=$(ps -o user= -p "$app_pid" | tr -d ' ')
if [ "$app_user" = "root" ]; then
  fail "the application is running as root — it must stay non-root"
fi

echo "PASS: config permissions are least-privilege and the app still runs as $app_user"
