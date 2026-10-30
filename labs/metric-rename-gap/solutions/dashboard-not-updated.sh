#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point the panel at the metric's current name
sed -i 's/^metric: .*/metric: http_request_count/' /etc/dashboards/overview.yaml
render-dashboard
