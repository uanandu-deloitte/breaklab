#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the hook's shebang and run the deploy it missed
sed -i '1s#.*#\#!/usr/bin/env bash#' /srv/deploy.git/hooks/post-receive
bash /srv/deploy.git/hooks/post-receive
