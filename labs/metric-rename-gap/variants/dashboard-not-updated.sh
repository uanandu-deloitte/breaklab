#!/usr/bin/env bash
set -euo pipefail
# Oops: the exporter's metric was renamed a while back; this dashboard's
# panel config was never updated to match
sed -i 's/^metric: .*/metric: http_requests_total/' /etc/dashboards/overview.yaml
rm -f /etc/dashboards/aliases.yaml
render-dashboard
