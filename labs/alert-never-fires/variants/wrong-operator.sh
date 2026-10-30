#!/usr/bin/env bash
set -euo pipefail
# Oops: a "simplify the rule" edit flipped the comparison direction
sed -i 's/^op: .*/op: "<"/' /etc/alerts/rules.yaml
