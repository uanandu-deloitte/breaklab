#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore read access to the homepage file
chmod 644 /var/www/html/index.html
