#!/usr/bin/env bash
set -euo pipefail
# Oops: someone piped the app's output through tee "for visibility" --
# a pipeline always needs a shell process to hold it together, so that
# shell (not the app) is what ends up as PID 1's direct child, and it
# never passes SIGTERM down to the app on the other end of the pipe.
printf 'python3 /usr/local/bin/app.py 2>&1 | tee -a /var/log/app/app.log\n' > /etc/app/launch-cmd
kill -HUP 1
sleep 0.5
