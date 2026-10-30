#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point current at the release that actually exists
ln -sfn /opt/releases/v2 /opt/current
