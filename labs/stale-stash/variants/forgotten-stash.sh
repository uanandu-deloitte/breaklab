#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Oops: the hotfix was stashed ("let me pull latest first") and never popped
sed -i 's/^timeout_seconds:.*/timeout_seconds: 60/' CONFIG.md
git stash push -q -m "hotfix: bump timeout_seconds to 60"
