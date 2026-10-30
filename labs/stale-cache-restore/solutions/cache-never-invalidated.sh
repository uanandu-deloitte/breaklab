#!/usr/bin/env bash
set -euo pipefail
# Reference fix: refresh the dependency cache to match the current lockfile,
# then rebuild vendor from it
cp /srv/app/requirements.lock /var/cache/ci/deps/requirements.lock
ci-build
