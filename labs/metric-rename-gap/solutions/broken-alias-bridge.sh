#!/usr/bin/env bash
set -euo pipefail
# Reference fix: correct the alias key so it actually matches what the
# dashboard declares
sed -i 's/^http_requests_totals:/http_requests_total:/' /etc/dashboards/aliases.yaml
render-dashboard
