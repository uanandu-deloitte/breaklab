#!/usr/bin/env bash
set -euo pipefail
# Oops: a verbose logging flag left on overnight filled the app log
dd if=/dev/zero of=/var/app/data/app.log bs=1M count=80 status=none
