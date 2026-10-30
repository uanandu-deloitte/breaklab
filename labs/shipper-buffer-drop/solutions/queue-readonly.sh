#!/usr/bin/env bash
set -euo pipefail
# Reference fix: clear the read-only marker
rm -f /var/spool/shipper/.readonly
