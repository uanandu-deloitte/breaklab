#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Reference fix: find the commit that broke the check and revert just that one
sha=$(git log --oneline --grep="Tighten healthcheck strictness" --format='%H')
git revert --no-edit "$sha" >/dev/null
