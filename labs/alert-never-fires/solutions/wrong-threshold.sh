#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore a sane threshold
sed -i 's/^threshold: .*/threshold: 90/' /etc/alerts/rules.yaml
alert-evaluator
