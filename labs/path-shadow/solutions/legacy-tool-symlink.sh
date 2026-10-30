#!/usr/bin/env bash
set -euo pipefail
# Reference fix: remove the shadowing symlink so the real jq on $PATH wins.
rm -f /usr/local/bin/jq
deploy-sync
