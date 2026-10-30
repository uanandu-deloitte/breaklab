#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Reference fix: find the orphaned commit via the HEAD reflog and merge it onto main
sha=$(git reflog --format='%H %gs' | grep "RELEASE: v1.2.0" | head -n1 | cut -d' ' -f1)
git merge -q --ff-only "$sha"
