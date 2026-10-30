#!/usr/bin/env bash
set -euo pipefail
# Oops: someone hand-edited the deployed file instead of going through infra-apply
echo "quick hotfix, will clean up later" > /var/www/html/index.html
nginx -t
nginx
