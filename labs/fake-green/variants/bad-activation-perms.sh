#!/usr/bin/env bash
set -euo pipefail
# Oops: a "lockdown" pass restricted /srv/app to root only
chown root:root /srv/app && chmod 755 /srv/app
su - deployer -c "cd ~/app-src && git push -q origin main"
