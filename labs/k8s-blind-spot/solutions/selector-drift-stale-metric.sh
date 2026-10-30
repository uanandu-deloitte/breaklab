#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the Service's real selector
kubectl patch service checkout -p '{"spec":{"selector":{"app":"checkout"}}}'
# Reference fix: point the refresh wrapper back at the real probe and run it
sed -i 's#/usr/local/bin/metrics-agent-legacy#/usr/local/bin/metrics-agent#' /usr/local/bin/refresh-metrics
refresh-metrics
