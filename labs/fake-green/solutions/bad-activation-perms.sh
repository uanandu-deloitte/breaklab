#!/usr/bin/env bash
set -euo pipefail
# Reference fix: give deployer write access back on the activation directory
chown deployer:deployer /srv/app
su - deployer -c "bash /srv/deploy.git/hooks/post-receive"
