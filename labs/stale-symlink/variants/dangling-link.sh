#!/usr/bin/env bash
set -euo pipefail
# Oops: the old release directory was cleaned up, leaving current dangling
rm -rf /opt/releases/v1
ln -sfn /opt/releases/v1 /opt/current
