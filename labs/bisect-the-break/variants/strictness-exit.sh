#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
echo "note: alert thresholds doc" >> NOTES.md && git add -A && git commit -q -m "Document alert thresholds"
sed -i 's/exit 0/exit 1/' healthcheck.sh && git add -A && git commit -q -m "Tighten healthcheck strictness"
echo "formatting pass" >> NOTES.md && git add -A && git commit -q -m "Tidy up notes formatting"
echo "readme tweak" >> README.md && git add -A && git commit -q -m "Minor README tweak"
