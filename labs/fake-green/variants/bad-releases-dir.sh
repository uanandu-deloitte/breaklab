#!/usr/bin/env bash
set -euo pipefail
# Oops: a path cleanup pointed releases at a directory deployer can't write to
sed -i 's#RELEASES_DIR=.*#RELEASES_DIR=/root/releases#' /etc/deploy/deploy.conf
su - deployer -c "cd ~/app-src && git push -q origin main"
