#!/usr/bin/env bash
set -euo pipefail
cd /srv/app-repo

# Find the hotfix wherever it actually landed, by content/message, and
# bring it onto release properly instead of guessing a SHA.
sha=$(git log develop --oneline --grep='^HOTFIX' | awk '{print $1}' | head -n1)
git checkout -q release
git cherry-pick "$sha" >/dev/null

# Fix the pipeline's cache key so it's derived from the real application
# source again, not an unrelated file that never changes.
sed -i 's#^HASH_SOURCE=.*#HASH_SOURCE=src/app.sh#' /opt/pipeline/build.sh

/opt/pipeline/build.sh
