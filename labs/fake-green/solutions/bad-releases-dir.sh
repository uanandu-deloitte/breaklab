#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point releases back at a directory deployer can actually write to
sed -i 's#RELEASES_DIR=.*#RELEASES_DIR=/srv/releases#' /etc/deploy/deploy.conf
su - deployer -c "bash /srv/deploy.git/hooks/post-receive"
