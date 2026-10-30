#!/usr/bin/env bash
set -euo pipefail
# Loads whatever environment this service is configured to run with,
# then writes its latest export for other processes to pick up.
if [ -f /etc/default/producer ]; then
  . /etc/default/producer
fi
for f in /etc/profile.d/*.sh; do
  [ -f "$f" ] && . "$f"
done
echo '{"ok":true}' > /var/lib/appdata/export.json
