#!/usr/bin/env bash
set -euo pipefail
# Reference fix: force the cache to be treated as stale so it actually gets
# refreshed from the current lockfile instead of trusting an unrelated mtime
touch /var/cache/ci/deps
ci-build
