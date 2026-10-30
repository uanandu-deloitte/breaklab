#!/usr/bin/env bash
set -euo pipefail
# Oops: the release script pointed current at the old release by mistake
ln -sfn /opt/releases/v1 /opt/current
