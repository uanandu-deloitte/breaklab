#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the hook's execute bit and run the deploy it missed
chmod 755 /srv/deploy.git/hooks/post-receive
bash /srv/deploy.git/hooks/post-receive
