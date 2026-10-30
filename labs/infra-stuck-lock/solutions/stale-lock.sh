#!/usr/bin/env bash
set -euo pipefail
rm -f /var/run/infra/apply.lock
infra-apply
