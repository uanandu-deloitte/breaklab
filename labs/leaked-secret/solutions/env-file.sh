#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Reference fix: scrub the secret file from every commit, not just delete it
# again. git-filter-repo is the Git project's own recommended replacement
# for filter-branch (faster, safer defaults, no manual ref cleanup needed).
git filter-repo --force --invert-paths --path config/secrets.env >/dev/null 2>&1
