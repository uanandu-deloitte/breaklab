#!/usr/bin/env bash
set -euo pipefail
echo "v2" > /infra/MODULE_VERSION
infra-apply
