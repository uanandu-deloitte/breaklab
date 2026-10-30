#!/usr/bin/env bash
set -euo pipefail
# Oops: a config edit pointed the shipper at a path that doesn't exist
sed -i 's#DEST=.*#DEST=/var/log/aggregator/shipped.log#' /etc/log-shipper.conf
