#!/usr/bin/env bash
set -euo pipefail
# Reference fix: deploy whatever is actually tagged as the current prod
# build, not whatever filename happens to sort last
latest=$(cat /srv/artifacts/PROD_LATEST)
rm -rf /srv/app/live
mkdir -p /srv/app/live
tar -xf "/srv/artifacts/$latest" -C /srv/app/live
echo "promote: deployed $latest to /srv/app/live"
