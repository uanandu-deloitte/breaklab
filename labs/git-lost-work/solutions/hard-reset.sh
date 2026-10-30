#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Reference fix: find the discarded commit in the reflog and fast-forward to it
sha=$(git reflog --format='%H %gs' | grep "RELEASE: v1.2.0" | head -n1 | cut -d' ' -f1)
git reset -q --hard "$sha"
