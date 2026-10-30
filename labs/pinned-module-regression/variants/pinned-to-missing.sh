#!/usr/bin/env bash
set -euo pipefail
# Oops: pinned at a module version that was since removed
# First, make the on-disk deployed config stale relative to a changed
# source, same as any real incident:
sed -i 's/^message:.*/message: Order confirmed and shipped today/' /infra/site.yaml
echo "v0" > /infra/MODULE_VERSION
infra-apply
