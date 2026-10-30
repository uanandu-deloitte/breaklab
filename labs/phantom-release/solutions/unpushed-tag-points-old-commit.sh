#!/usr/bin/env bash
set -euo pipefail
cd /srv/app-repo

# Move the tag the pipeline actually builds from onto release's real tip,
# which already has the hotfix -- nothing else depends on the old
# v2.0.0 yet, so re-pointing it is safe here.
git tag -f v2.0.0 release

/opt/pipeline/build.sh
