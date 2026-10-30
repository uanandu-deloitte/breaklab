#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
# Oops: the release was committed while HEAD was detached, never on a branch
git checkout -q --detach main
echo "v1.2.0" > RELEASE.md
git add -A
git commit -q -m "RELEASE: v1.2.0 going out tonight"
git checkout -q main
