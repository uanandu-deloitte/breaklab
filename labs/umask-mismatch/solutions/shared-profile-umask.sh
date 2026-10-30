#!/usr/bin/env bash
set -euo pipefail
# Reference fix: that hardening snippet shouldn't apply to this
# service's export at all.
rm -f /etc/profile.d/10-harden.sh
