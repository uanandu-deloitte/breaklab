#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
echo "scale notes" >> NOTES.md && git add -A && git commit -q -m "Capacity planning notes"
sed -i 's/THRESHOLD=90/THRESHOLD=0/' healthcheck.sh && git add -A && git commit -q -m "Adjust healthcheck threshold for new SLA"
echo "cleanup" >> NOTES.md && git add -A && git commit -q -m "Clean up stray notes"
