#!/usr/bin/env bash
set -euo pipefail
# Oops: pinned at a module version that was since removed. First, make
# the on-disk deployed config stale relative to a changed source, same
# as any real incident:
sed -i 's/^message:.*/message: Winter sale: 30% off everything/' /infra/site.yaml
echo "v0" > /infra/MODULE_VERSION
infra-apply

cd /srv/release-repo
# Oops: the release was committed while HEAD was detached, never on a branch
git checkout -q --detach main
echo "v1.2.0" > RELEASE_MESSAGE.md
git add -A
git commit -q -m "RELEASE: v1.2.0 going out tonight"
git checkout -q main
