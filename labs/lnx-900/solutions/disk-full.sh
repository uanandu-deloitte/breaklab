#!/usr/bin/env bash
set -euo pipefail
# Reference fix: remove only the disposable junk, leave retained data alone.
rm -rf /var/log/app-data/junk
curl -fsS http://127.0.0.1:5000/ >/dev/null
