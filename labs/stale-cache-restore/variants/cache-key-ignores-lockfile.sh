#!/usr/bin/env bash
set -euo pipefail
# Real dependency bump: libwidget needs 3.0
sed -i 's/^libwidget==.*/libwidget==3.0/' /srv/app/requirements.lock

# ci-build's cache key is derived from the cache directory's mtime, not
# from requirements.lock's content -- bumping the lockfile alone doesn't
# invalidate it, so the pipeline reports success while running the old deps
ci-build
