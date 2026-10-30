#!/usr/bin/env bash
set -euo pipefail
# Oops: a cleanup pass removed what looked like an unused package
mv /usr/bin/jq /usr/bin/jq.disabled
