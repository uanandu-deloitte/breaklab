#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Reference fix: scrub the hardcoded key from every commit, not just the
# latest one. git-filter-repo is the Git project's own recommended
# replacement for filter-branch's --tree-filter (a single pass over every
# blob instead of a full checkout-and-rerun per commit).
printf 'sk-FAKE9876543210==>REDACTED\n' > /tmp/replace.txt
git filter-repo --force --replace-text /tmp/replace.txt >/dev/null 2>&1
