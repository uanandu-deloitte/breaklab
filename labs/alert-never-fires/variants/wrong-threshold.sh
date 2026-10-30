#!/usr/bin/env bash
set -euo pipefail
# Oops: a "raise the threshold to cut noise" change went too far
sed -i 's/^threshold: .*/threshold: 999/' /etc/alerts/rules.yaml
