#!/usr/bin/env bash
set -euo pipefail
job="$1"
mv "$job" "/var/spool/jobs/done/$(basename "$job")"
