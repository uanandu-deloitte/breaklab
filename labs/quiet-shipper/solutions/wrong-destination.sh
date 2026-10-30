#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point the shipper back at the real destination and re-run it
sed -i 's#DEST=.*#DEST=/var/log/shipped.log#' /etc/log-shipper.conf
ship-logs
