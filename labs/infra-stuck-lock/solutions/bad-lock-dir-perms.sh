#!/usr/bin/env bash
set -euo pipefail
chown infraops:infraops /var/run/infra
chmod 755 /var/run/infra
infra-apply
