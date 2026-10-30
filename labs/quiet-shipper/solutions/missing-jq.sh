#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore jq and re-run the shipper
mv /usr/bin/jq.disabled /usr/bin/jq
ship-logs
