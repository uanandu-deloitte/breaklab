#!/usr/bin/env bash
set -euo pipefail
# This quarter's message changed
sed -i 's/^message:.*/message: Q3 pricing notice v2/' /infra/site.yaml
# Simulate the usual deploy habit: everyone runs apply-site, not
# infra-apply directly
apply-site
