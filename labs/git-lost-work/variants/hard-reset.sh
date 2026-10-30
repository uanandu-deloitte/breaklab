#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
echo "v1.2.0" > RELEASE.md
git add -A
git commit -q -m "RELEASE: v1.2.0 going out tonight"
# Oops: an accidental hard reset threw the release commit away
git reset -q --hard HEAD~1
