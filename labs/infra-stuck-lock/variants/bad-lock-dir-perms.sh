#!/usr/bin/env bash
set -euo pipefail
# Oops: a "harden permissions" pass locked down /var/run/infra too far
chown root:root /var/run/infra
chmod 700 /var/run/infra
# Today's config changed; the next apply should have picked it up
sed -i 's/^message:.*/message: BreakLab OK v5/' /infra/site.yaml
