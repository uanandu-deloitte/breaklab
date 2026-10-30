#!/usr/bin/env bash
set -euo pipefail
# Oops: the source was updated for tonight's change, but nobody re-applied it
sed -i 's/message: .*/message: BreakLab OK v4/' /infra/site.yaml
nginx -t
nginx
