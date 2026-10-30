#!/usr/bin/env bash
set -euo pipefail
echo "v2" > /infra/MODULE_VERSION
infra-apply

cd /srv/release-repo
sha=$(git reflog --format='%H %gs' | grep "RELEASE: v1.2.0" | head -n1 | cut -d' ' -f1)
git merge -q --ff-only "$sha"
