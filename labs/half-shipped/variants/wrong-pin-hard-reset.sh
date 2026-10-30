#!/usr/bin/env bash
set -euo pipefail
# Oops: pinned back to the old render module
echo "v1" > /infra/MODULE_VERSION
infra-apply

cd /srv/release-repo
echo "v1.2.0" > RELEASE_MESSAGE.md
git add -A
git commit -q -m "RELEASE: v1.2.0 going out tonight"
# Oops: an accidental hard reset threw the release commit away
git reset -q --hard HEAD~1
