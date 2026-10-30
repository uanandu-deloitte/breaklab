#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
mkdir -p config
printf 'API_KEY=sk-FAKE1234567890\n' > config/secrets.env
git add -A && git commit -q -m "Add service config"
echo "note" >> NOTES.md && git add -A && git commit -q -m "Unrelated housekeeping"
git rm -q config/secrets.env
git commit -q -m "Remove secrets file from repo"
