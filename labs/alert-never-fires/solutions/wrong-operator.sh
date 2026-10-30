#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the correct comparison direction
sed -i 's/^op: .*/op: ">"/' /etc/alerts/rules.yaml
alert-evaluator
