#!/usr/bin/env bash
set -euo pipefail
# Oops: a permissions "hardening" pass locked deployer out of its own lock directory
chown root:root /var/run/deploy
chmod 700 /var/run/deploy
