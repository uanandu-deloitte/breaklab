#!/usr/bin/env bash
set -euo pipefail
sed -i 's/^SITE_MESSAGE="[^"]*" //' /usr/local/bin/apply-site
infra-apply
