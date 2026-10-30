#!/usr/bin/env bash
set -euo pipefail
# Oops: pinned back to the old render module
echo "v1" > /infra/MODULE_VERSION
infra-apply
