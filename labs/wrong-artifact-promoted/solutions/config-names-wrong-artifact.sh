#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point the promote config back at the artifact that's
# actually tagged prod
latest=$(cat /srv/artifacts/PROD_LATEST)
sed -i "s/^PROD_ARTIFACT=.*/PROD_ARTIFACT=$latest/" /etc/deploy/promote.conf
promote
