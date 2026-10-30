#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the deploy lock directory to its owning user
chown deployer:deployer /var/run/deploy
chmod 755 /var/run/deploy
deploy
