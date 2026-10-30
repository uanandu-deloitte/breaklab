#!/usr/bin/env bash
set -euo pipefail
# Oops: an edit to the hook corrupted its shebang line
sed -i '1s#.*#\#!/usr/bin/env not-a-real-interpreter#' /srv/deploy.git/hooks/post-receive
# The engineer's push from earlier tonight, which the broken hook couldn't run
cd /root/app-src && git push -q origin main
