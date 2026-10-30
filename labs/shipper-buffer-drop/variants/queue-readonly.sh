#!/usr/bin/env bash
set -euo pipefail
# Oops: an automated hardening pass marked the spool read-only, and
# ship-logs fails open instead of erroring loudly
touch /var/spool/shipper/.readonly
