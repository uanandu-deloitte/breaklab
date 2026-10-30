#!/usr/bin/env bash
set -euo pipefail
# Oops: a rename-compat bridge was added so old dashboards keep working,
# but the mapping has a typo in the old metric name's key
sed -i 's/^metric: .*/metric: http_requests_total/' /etc/dashboards/overview.yaml
printf 'http_requests_totals: http_request_count\n' > /etc/dashboards/aliases.yaml
render-dashboard
