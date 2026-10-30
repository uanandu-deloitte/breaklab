#!/usr/bin/env bash
set -euo pipefail
# Oops: a previous apply crashed mid-run and left its lock behind
su infraops -s /bin/bash -c "touch /var/run/infra/apply.lock"
# Today's config changed; the next apply should have picked it up
sed -i 's/^message:.*/message: BreakLab OK v5/' /infra/site.yaml
