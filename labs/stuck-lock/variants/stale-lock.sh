#!/usr/bin/env bash
set -euo pipefail
# Oops: the previous deploy crashed and never cleaned up its lock
su -s /bin/sh deployer -c 'touch /var/run/deploy/deploy.lock'
