#!/usr/bin/env bash
set -euo pipefail
# Reference fix: the previous deploy is long gone; clear its stale lock and redeploy
rm -f /var/run/deploy/deploy.lock
deploy
